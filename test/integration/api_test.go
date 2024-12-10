package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rossi1/ensync-cli/internal/api"
	"github.com/rossi1/ensync-cli/internal/domain"
)

type mockHandler func(t *testing.T, w http.ResponseWriter, r *http.Request)

func TestClient(t *testing.T) {
	mockServer := setupMockServer(t)
	defer mockServer.Close()

	client := api.NewClient(mockServer.URL, "test-api-key")
	ctx := context.Background()

	t.Run("Events", func(t *testing.T) {
		t.Run("List", testListEvents(ctx, client))
		t.Run("Create", testCreateEvent(ctx, client))
		t.Run("Update", testUpdateEvent(ctx, client))

		t.Run("Update", testUpdateEvent(ctx, client))

		t.Run("Name", testGetEventByName(ctx, client))
	})

	t.Run("AccessKeys", func(t *testing.T) {
		t.Run("List", testListAccessKeys(ctx, client))
		t.Run("Create", testCreateAccessKey(ctx, client))
		t.Run("GetPermissions", testGetAccessKeyPermissions(ctx, client))
		t.Run("SetPermissions", testSetAccessKeyPermissions(ctx, client))
	})

	t.Run("ErrorCases", testErrorCases)
}

func setupMockServer(t *testing.T) *httptest.Server {
	handlers := map[string]mockHandler{
		"GET/event":                   mockListEvents,
		"POST/event":                  mockCreateEvent,
		"PUT/event":                   mockUpdateEvent,
		"GET/event/name":              mockGetEventByName,
		"GET/access-key":              mockListAccessKeys,
		"POST/access-key":             mockCreateAccessKey,
		"GET/access-key/permissions":  mockGetAccessKeyPermissions,
		"POST/access-key/permissions": mockSetAccessKeyPermissions,
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		key := r.Method + r.URL.Path
		for pattern, handler := range handlers {
			if strings.HasPrefix(key, pattern) {
				handler(t, w, r)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	}))
}

func testListEvents(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		params := &api.ListParams{
			PageIndex: 0,
			Limit:     10,
			Order:     "DESC",
			OrderBy:   "createdAt",
		}

		events, err := client.ListEvents(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, events)
		assert.Equal(t, 2, events.ResultsLength)
		assert.Len(t, events.Results, 2)
	}
}

func testGetEventByName(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		name := "test-event"
		event, err := client.GetEventByName(ctx, name)
		require.NoError(t, err)
		require.NotNil(t, event)
		assert.Equal(t, name, event.Name)
	}
}

func testCreateEvent(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		event := &domain.Event{
			Name: "test-event",
			Payload: map[string]string{
				"key": "value",
			},
		}

		err := client.CreateEvent(ctx, event)
		require.NoError(t, err)
	}
}

func testUpdateEvent(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		event := &domain.Event{
			ID:   123,
			Name: "updated-event",
			Payload: map[string]string{
				"key": "new-value",
			},
		}

		err := client.UpdateEvent(ctx, event)
		require.NoError(t, err)
	}
}

func testListAccessKeys(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		params := &api.ListParams{
			PageIndex: 0,
			Limit:     10,
			Order:     "DESC",
			OrderBy:   "createdAt",
		}

		keys, err := client.ListAccessKeys(ctx, params)
		require.NoError(t, err)
		require.NotNil(t, keys)
		assert.Equal(t, 2, keys.ResultsLength)
	}
}

func testCreateAccessKey(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		permissions := &domain.AccessKeyPermissions{
			Send:    []string{"event1", "event2"},
			Receive: []string{"event3"},
		}
		key := &domain.AccessKey{Permissions: permissions}

		created, err := client.CreateAccessKey(ctx, key)
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.Key)
	}
}

func testGetAccessKeyPermissions(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		permissions, err := client.GetAccessKeyPermissions(ctx, "test-key")
		require.NoError(t, err)
		require.NotNil(t, permissions)
		assert.Contains(t, permissions.Send, "event1")
	}
}

func testSetAccessKeyPermissions(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		permissions := &domain.AccessKeyPermissions{
			Send:    []string{"event1", "event2"},
			Receive: []string{"event3"},
		}
		err := client.SetAccessKeyPermissions(ctx, "test-key", permissions)
		require.NoError(t, err)
	}
}

func mockListEvents(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "GET", "test-api-key")
	sendJSONResponse(w, domain.EventList{
		ResultsLength: 2,
		Results: []*domain.Event{
			{
				ID:      1,
				Name:    "test-event-1",
				Payload: map[string]string{"key": "value1"},
			},
			{
				ID:      2,
				Name:    "test-event-2",
				Payload: map[string]string{"key": "value2"},
			},
		},
	})
}

func mockCreateEvent(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "POST", "test-api-key")
	w.WriteHeader(http.StatusOK)
}

func mockUpdateEvent(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "PUT", "test-api-key")
	w.WriteHeader(http.StatusOK)
}

func mockListAccessKeys(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "GET", "test-api-key")
	sendJSONResponse(w, domain.AccessKeyList{
		ResultsLength: 2,
		Results: []*domain.AccessKey{
			{
				Key: "key1",
				Permissions: &domain.AccessKeyPermissions{
					Send:    []string{"event1"},
					Receive: []string{"event2"},
				},
			},
			{
				Key: "key2",
				Permissions: &domain.AccessKeyPermissions{
					Send:    []string{"event3"},
					Receive: []string{"event4"},
				},
			},
		},
	})
}

func mockGetEventByName(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "GET", "test-api-key")
	assert.Equal(t, "test-event", r.URL.Query().Get("name"))

	createdAt, err := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("Failed to parse createdAt: %v", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	if err != nil {
		t.Fatalf("Failed to parse updatedAt: %v", err)
	}

	sendJSONResponse(w, domain.Event{
		ID:        1,
		Name:      "test-event",
		Payload:   map[string]string{"key": "value"},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
}

func mockCreateAccessKey(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "POST", "test-api-key")
	sendJSONResponse(w, map[string]string{"accessKey": "new-access-key-123"})
}

func mockGetAccessKeyPermissions(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "GET", "test-api-key")
	sendJSONResponse(w, map[string]interface{}{
		"permissions": &domain.AccessKeyPermissions{
			Send:    []string{"event1", "event2"},
			Receive: []string{"event3"},
		},
		"key":       "test-key",
		"createdAt": "2024-01-01T00:00:00Z",
	})
}

func mockSetAccessKeyPermissions(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "POST", "test-api-key")
	w.WriteHeader(http.StatusOK)
}

func verifyRequest(t *testing.T, r *http.Request, method, expectedAPIKey string) {
	assert.Equal(t, method, r.Method)
	assert.Equal(t, expectedAPIKey, r.Header.Get("X-API-KEY"))
}

func sendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func testErrorCases(t *testing.T) {
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid API key",
		})
	}))
	defer errorServer.Close()

	client := api.NewClient(errorServer.URL, "invalid-key")
	ctx := context.Background()

	t.Run("ListEvents", func(t *testing.T) {
		events, err := client.ListEvents(ctx, &api.ListParams{})
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Contains(t, err.Error(), "Invalid API key")
	})

	t.Run("ListAccessKeys", func(t *testing.T) {
		keys, err := client.ListAccessKeys(ctx, &api.ListParams{})
		assert.Error(t, err)
		assert.Nil(t, keys)
		assert.Contains(t, err.Error(), "Invalid API key")
	})
}
