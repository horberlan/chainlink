package periodicbackup

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var (
	fileName           = "db_backup.tar.gz"
	minBackupFrequency = time.Minute
)

type backupResult struct {
	size int64
	path string
}

type (
	PeriodicBackup interface {
		Start() error
		Close() error
	}

	periodicBackup struct {
		logger          *logger.Logger
		databaseURL     url.URL
		frequency       time.Duration
		outputParentDir string
		done            chan bool
	}
)

func NewPeriodicBackup(frequency time.Duration, databaseURL url.URL, outputParentDir string, logger *logger.Logger) PeriodicBackup {
	return &periodicBackup{
		logger,
		databaseURL,
		frequency,
		outputParentDir,
		make(chan bool),
	}
}

func (backup periodicBackup) Start() error {
	backup.runBackupGracefully()

	ticker := time.NewTicker(backup.frequency)

	go func() {
		for {
			select {
			case <-backup.done:
				ticker.Stop()
				return
			case <-ticker.C:
				if backup.frequencyIsTooSmall() {
					logger.Errorf("Database backup frequency (%s=%v) is too small. Please set it to at least %s", "DATABASE_BACKUP_FREQUENCY", backup.frequency, minBackupFrequency)
					continue
				}
				backup.runBackupGracefully()
			}
		}
	}()

	return nil
}

func (backup periodicBackup) Close() error {
	backup.done <- true
	return nil
}

func (backup *periodicBackup) frequencyIsTooSmall() bool {
	return backup.frequency < minBackupFrequency
}

func (backup *periodicBackup) runBackupGracefully() {
	backup.logger.Info("PeriodicBackup: Starting database backup...")
	startAt := time.Now()
	result, err := backup.runBackup()
	duration := time.Since(startAt)
	if err != nil {
		backup.logger.Errorf("PeriodicBackup: Failed after %s with: %v", duration, err)
	} else {
		backup.logger.Infof("PeriodicBackup: Database backup finished successfully after %s: %d bytes written to %s", duration, result.size, result.path)
		if duration > backup.frequency {
			backup.logger.Warn("PeriodicBackup: Backup is taking longer to complete than the frequency")
		}
	}
}

func (backup *periodicBackup) runBackup() (*backupResult, error) {

	tmpFile, err := ioutil.TempFile(backup.outputParentDir, "db_backup")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a tmp file")
	}
	err = os.Remove(tmpFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to remove the tmp file before running backup")
	}

	cmd := exec.Command(
		"pg_dump", backup.databaseURL.String(),
		"-f", tmpFile.Name(),
		"-F", "t", // format: tar
	)

	_, err = cmd.Output()

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, errors.Wrap(err, fmt.Sprintf("pg_dump failed with output: %s", string(ee.Stderr)))
		}
		return nil, errors.Wrap(err, "pg_dump failed")
	}

	finalFilePath := filepath.Join(backup.outputParentDir, fileName)
	_ = os.Remove(finalFilePath)
	err = os.Rename(tmpFile.Name(), finalFilePath)
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return nil, errors.Wrap(err, "Failed to rename the temp file to the final backup file")
	}

	file, err := os.Stat(finalFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to access the final backup file")
	}

	return &backupResult{
		size: file.Size(),
		path: finalFilePath,
	}, nil
}
