package periodicbackup

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriodicBackup_RunBackup(t *testing.T) {
	rawConfig := orm.NewConfig()
	periodicBackup := NewPeriodicBackup(time.Minute, rawConfig.DatabaseURL(), os.TempDir(), logger.Default).(*periodicBackup)
	assert.False(t, periodicBackup.frequencyIsTooSmall())

	result, err := periodicBackup.runBackup()
	require.NoError(t, err, "error not nil for backup")

	defer os.Remove(result.path)

	file, err := os.Stat(result.path)
	require.NoError(t, err, "error not nil when checking for output file")

	assert.Greater(t, file.Size(), int64(0))
	assert.Equal(t, file.Size(), result.size)
	assert.True(t, strings.Contains(result.path, fileName))
}

func TestPeriodicBackup_FrequencyTooSmall(t *testing.T) {
	rawConfig := orm.NewConfig()
	periodicBackup := NewPeriodicBackup(time.Second, rawConfig.DatabaseURL(), os.TempDir(), logger.Default).(*periodicBackup)
	assert.True(t, periodicBackup.frequencyIsTooSmall())
}
