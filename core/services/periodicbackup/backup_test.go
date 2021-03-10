package periodicbackup

import (
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"gotest.tools/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)


func TestPeriodicBackup_Run_Backup(t *testing.T) {
	rawConfig := orm.NewConfig()
	periodicBackup := NewBackgroundBackup(time.Minute, rawConfig.DatabaseURL(), os.TempDir(), logger.Default)
	err := periodicBackup.RunBackup()

	defer os.Remove(filepath.Join(os.TempDir(), fileName))
	assert.NilError(t, err, "error not nil")
}
