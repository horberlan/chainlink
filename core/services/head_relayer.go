package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const callbackTimeout = 2 * time.Second

type callbackID [256]byte

// HeadRelayable defines the interface for listeners
type HeadRelayable interface {
	OnNewLongestChain(ctx context.Context, head models.Head)
}

type callbackSet map[callbackID]HeadRelayable

// NewHeadRelayer creates a new HeadRelayer
func NewHeadRelayer() *HeadRelayer {
	return &HeadRelayer{
		callbacks:     make(callbackSet),
		mailbox:       utils.NewMailbox(1),
		mutex:         &sync.RWMutex{},
		chClose:       make(chan struct{}),
		wgDone:        sync.WaitGroup{},
		StartStopOnce: utils.StartStopOnce{},
	}
}

// HeadRelayer relays heads from the head tracker to subscribed jobs, it is less robust against
// congestion than the head tracker, and missed heads should be expected by consuming jobs
type HeadRelayer struct {
	callbacks callbackSet
	mailbox   *utils.Mailbox
	mutex     *sync.RWMutex
	chClose   chan struct{}
	wgDone    sync.WaitGroup
	utils.StartStopOnce
}

var _ store.HeadTrackable = (*HeadRelayer)(nil)

func (hr *HeadRelayer) Start() error {
	return hr.StartOnce("HeadRelayer", func() error {
		hr.wgDone.Add(1)
		go hr.run()
		return nil
	})
}

func (hr *HeadRelayer) Close() error {
	if !hr.OkayToStop() {
		return errors.New("HeadRelayer is already stopped")
	}
	close(hr.chClose)
	hr.wgDone.Wait()
	return nil
}

func (hr *HeadRelayer) Connect(head *models.Head) error {
	return nil
}

func (hr *HeadRelayer) Disconnect() {}

func (hr *HeadRelayer) OnNewLongestChain(ctx context.Context, head models.Head) {
	hr.mailbox.Deliver(head)
}

func (hr *HeadRelayer) Subscribe(callback HeadRelayable) (unsubscribe func()) {
	hr.mutex.Lock()
	defer hr.mutex.Unlock()
	id, err := newID()
	if err != nil {
		logger.Errorf("Unable to create ID for head relayble callback: %v", err)
		return
	}
	hr.callbacks[id] = callback
	return func() {
		hr.mutex.Lock()
		defer hr.mutex.Unlock()
		delete(hr.callbacks, id)
	}
}

func (hr *HeadRelayer) run() {
	defer hr.wgDone.Done()
	for {
		select {
		case <-hr.chClose:
			return
		case <-hr.mailbox.Notify():
			hr.executeCallbacks()
		}
	}
}

// DEV: the head relayer makes no promises about head delivery! Subscribing
// Jobs should expect to the relayer to skip heads if there is a large number of listeners
// and all callbacks cannot be completed in the allotted time.
func (hr *HeadRelayer) executeCallbacks() {
	hr.mutex.RLock()
	callbacks := hr.copyCallbacks()
	hr.mutex.RUnlock()

	head, ok := hr.mailbox.Retrieve().(models.Head)
	if !ok {
		logger.Errorf("expected `models.Head`, got %T", head)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(hr.callbacks))

	for _, callback := range callbacks {
		go func(hr HeadRelayable) {
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), callbackTimeout)
			defer cancel()
			hr.OnNewLongestChain(ctx, head)
			elapsed := time.Since(start)
			logger.Debugw(fmt.Sprintf("HeadRelayer: finished callback in %s", elapsed), "callbackType", reflect.TypeOf(hr), "blockNumber", head.Number, "time", elapsed, "id", "head_relayer")
			wg.Done()
		}(callback)
	}

	wg.Wait()
}

func (hr *HeadRelayer) copyCallbacks() callbackSet {
	cp := make(callbackSet)
	for id, callback := range hr.callbacks {
		cp[id] = callback
	}
	return cp
}

func newID() (id callbackID, _ error) {
	randBytes := make([]byte, 256)
	_, err := rand.Read(randBytes)
	if err != nil {
		return id, err
	}
	copy(id[:], randBytes)
	return id, nil
}
