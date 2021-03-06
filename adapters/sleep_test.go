package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestSleep_Perform(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	adapter := adapters.Sleep{}
	err := json.Unmarshal([]byte(`{"until": 872835240}`), &adapter)
	assert.NoError(t, err)

	result := adapter.Perform(models.RunResult{}, store)
	assert.Equal(t, string(models.RunStatusPendingSleep), string(result.Status))
}
