package periodicbackup

import (
  "fmt"
  "github.com/pkg/errors"
  "github.com/smartcontractkit/chainlink/core/logger"
  "io/ioutil"
  "net/url"
  "os"
  "os/exec"
  "path/filepath"
  "time"
)

var (
  fileName = "db_backup.tar.gz"
)

type PeriodicBackup struct {
  logger *logger.Logger
  databaseURL url.URL
  frequency time.Duration
  outputParentDir string
  done chan bool
}

func NewBackgroundBackup(frequency time.Duration, databaseURL url.URL, outputParentDir string, logger *logger.Logger) PeriodicBackup {
  if frequency < time.Minute {
    logger.Fatalf("Database backup setting (%s=%v) is too frequent. Please set it to at least one minute.", "DATABASE_BACKUP_FREQUENCY", frequency)
  }
  return PeriodicBackup {
    logger,
    databaseURL,
    frequency,
    outputParentDir,
    make(chan bool),
  }
}


func (backup PeriodicBackup) Start() error {
  backup.RunBackupGracefully()

  ticker := time.NewTicker(backup.frequency)

  go func() {
    for {
      select {
      case <-backup.done:
        ticker.Stop()
        return
      case <-ticker.C:
        backup.RunBackupGracefully()
      }
    }
  }()

  return nil
}

func (backup PeriodicBackup) Close() error {
  backup.done <- true
  return nil
  // what if backup is just running?
}

func (backup *PeriodicBackup) RunBackupGracefully() {
  backup.logger.Info("PeriodicBackup: Running database backup...")
  err := backup.RunBackup()
  if err != nil {
    backup.logger.Errorf("PeriodicBackup: Failed: %v", err)
  } else {
    backup.logger.Info("PeriodicBackup: Database backup finished successfully")
  }
}

func (backup *PeriodicBackup) RunBackup() error {

  tmpFile, err := ioutil.TempFile(backup.outputParentDir, "db_backup")
  if err != nil {
    return errors.Wrap(err, "Failed to create a tmp file")
  }
  err = os.Remove(tmpFile.Name())
  if err != nil {
    return errors.Wrap(err, "Failed to remove the tmp file before running backup")
  }

  cmd := exec.Command(
    "pg_dump", backup.databaseURL.String(),
    "-f", tmpFile.Name(),
    "-F", "t", // format: tar
  )

  _, err = cmd.Output()

  if err != nil {
    if ee, ok := err.(*exec.ExitError); ok {
      return errors.Wrap(err, fmt.Sprintf("pg_dump failed with output: %s", string(ee.Stderr)))
    }
    return errors.Wrap(err, "pg_dump failed")
  }

  defer os.Remove(tmpFile.Name())

  finalFilePath := filepath.Join(backup.outputParentDir, fileName)
  _ = os.Remove(finalFilePath)
  err = os.Rename(tmpFile.Name(), finalFilePath)
  if err != nil {
    return errors.Wrap(err, "Failed to rename the temp file to the final backup file")
  }

  //TODO: verify that the file exists
  return nil
}