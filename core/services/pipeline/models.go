package pipeline

import (
	"fmt"
	"strconv"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

type Spec struct {
	ID              int32           `gorm:"primary_key"`
	DotDagSource    string          `json:"dotDagSource"`
	CreatedAt       time.Time       `json:"-"`
	MaxTaskDuration models.Interval `json:"-"`
	//PipelineTaskSpecs []TaskSpec      `json:"-" gorm:"foreignkey:PipelineSpecID;->"`
}

func (Spec) TableName() string {
	return "pipeline_specs"
}

// DEPRECATED
type TaskSpec struct {
	ID             int32             `json:"-" gorm:"primary_key"`
	DotID          string            `json:"dotId"`
	PipelineSpecID int32             `json:"-"`
	PipelineSpec   Spec              `json:"-"`
	Type           TaskType          `json:"-"`
	JSON           JSONSerializable  `json:"-" gorm:"type:jsonb"`
	Index          int32             `json:"-"`
	SuccessorID    null.Int          `json:"-"`
	CreatedAt      time.Time         `json:"-"`
	BridgeName     *string           `json:"-"`
	Bridge         models.BridgeType `json:"-" gorm:"foreignKey:BridgeName;->"`
}

func (TaskSpec) TableName() string {
	return "pipeline_task_specs"
}

type Run struct {
	ID               int64            `json:"-" gorm:"primary_key"`
	PipelineSpecID   int32            `json:"-"`
	PipelineSpec     Spec             `json:"pipelineSpec"`
	Meta             JSONSerializable `json:"meta"`
	Errors           JSONSerializable `json:"errors" gorm:"type:jsonb"`
	Outputs          JSONSerializable `json:"outputs" gorm:"type:jsonb"`
	CreatedAt        time.Time        `json:"createdAt"`
	FinishedAt       *time.Time       `json:"finishedAt"`
	PipelineTaskRuns []TaskRun        `json:"taskRuns" gorm:"foreignkey:PipelineRunID;->"`
}

func (Run) TableName() string {
	return "pipeline_runs"
}

func (r Run) GetID() string {
	return fmt.Sprintf("%v", r.ID)
}

func (r *Run) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	r.ID = int64(ID)
	return nil
}

func (r Run) HasErrors() bool {
	return r.FinalErrors().HasErrors()
}

func (r Run) FinalErrors() (f FinalErrors) {
	f, _ = r.Errors.Val.(FinalErrors)
	return f
}

type TaskRun struct {
	ID            int64             `json:"-" gorm:"primary_key"`
	Type          TaskType          `json:"type"`
	PipelineRun   Run               `json:"-"`
	PipelineRunID int64             `json:"-"`
	Output        *JSONSerializable `json:"output" gorm:"type:jsonb"`
	Error         null.String       `json:"error"`
	CreatedAt     time.Time         `json:"createdAt"`
	FinishedAt    *time.Time        `json:"finishedAt"`
	Index         int32

	// New
	DotID string

	// Deprecated
	//PipelineTaskSpecID int32             `json:"-"`
	//PipelineTaskSpec   TaskSpec          `json:"taskSpec" gorm:"foreignkey:PipelineTaskSpecID;->"`
}

func (TaskRun) TableName() string {
	return "pipeline_task_runs"
}

// RunStatus represents the status of a run
type RunStatus int

const (
	// RunStatusUnknown is the when the run status cannot be determined.
	RunStatusUnknown RunStatus = iota
	// RunStatusInProgress is used for when a run is actively being executed.
	RunStatusInProgress
	// RunStatusErrored is used for when a run has errored and will not complete.
	RunStatusErrored
	// RunStatusCompleted is used for when a run has successfully completed execution.
	RunStatusCompleted
)

// Completed returns true if the status is RunStatusCompleted.
func (s RunStatus) Completed() bool {
	return s == RunStatusCompleted
}

// Errored returns true if the status is RunStatusErrored.
func (s RunStatus) Errored() bool {
	return s == RunStatusErrored
}

// Finished returns true if the status is final and can't be changed.
func (s RunStatus) Finished() bool {
	return s.Completed() || s.Errored()
}

// Status determines the status of the run.
func (r *Run) Status() RunStatus {
	if r.HasErrors() {
		return RunStatusErrored
	} else if r.FinishedAt != nil {
		return RunStatusCompleted
	}

	return RunStatusInProgress
}

func (tr TaskRun) GetID() string {
	return fmt.Sprintf("%v", tr.ID)
}

func (tr *TaskRun) SetID(value string) error {
	ID, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	tr.ID = int64(ID)
	return nil
}

//func (s TaskSpec) IsFinalPipelineOutput() bool {
//	return s.SuccessorID.IsZero()
//}

func (tr TaskRun) GetDotID() string {
	//return tr.PipelineTaskSpec.GetDotID
	return tr.DotID
}

func (tr TaskRun) Result() Result {
	var result Result
	if !tr.Error.IsZero() {
		result.Error = errors.New(tr.Error.ValueOrZero())
	} else if tr.Output != nil && tr.Output.Val != nil {
		result.Value = tr.Output.Val
	}
	return result
}
