package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ensync-cli/internal/api"
)

func TestClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := api.NewClient(
		"http://localhost:8080/api/v1/ensync",
		"test-api-key",
	)

	t.Run("ListEvents", func(t *testing.T) {
		params := &api.ListParams{
			PageIndex: 0,
			Limit:     10,
			Order:     "DESC",
			OrderBy:   "createdAt",
		}

		events, err := client.ListEvents(context.Background(), params)
		require.NoError(t, err)
		assert.NotNil(t, events)
		assert.GreaterOrEqual(t, events.ResultsLength, 0)
	})
}
