package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rossi1/ensync-cli/internal/api"
	"github.com/rossi1/ensync-cli/internal/domain"
)

type mockHandler func(t *testing.T, w http.ResponseWriter, r *http.Request)

func setupMockServer(t *testing.T) *httptest.Server {
	handlers := map[string]mockHandler{
		"GET/event":       mockListEvents,
		"POST/event":      mockCreateEvent,
		"GET/access-key":  mockListAccessKeys,
		"POST/access-key": mockCreateAccessKey,
	}

	dynamicHandlers := map[*regexp.Regexp]map[string]mockHandler{
		regexp.MustCompile(`^/access-key/permissions/[\w-]+$`): {
			http.MethodGet:  mockGetAccessKeyPermissions,
			http.MethodPost: mockSetAccessKeyPermissions,
		},
		regexp.MustCompile(`^/event/[\w-]+$`): {
			http.MethodPut: mockUpdateEvent,
			http.MethodGet: mockGetEventByName,
		},
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		key := r.Method + r.URL.Path
		if handler, exists := handlers[key]; exists {
			handler(t, w, r)
			return
		}

		for regex, methodHandlers := range dynamicHandlers {
			if regex.MatchString(r.URL.Path) {
				if handler, exists := methodHandlers[r.Method]; exists {
					handler(t, w, r)
					return
				}
			}
		}

	}))
}

func TestClient(t *testing.T) {
	mockServer := setupMockServer(t)
	defer mockServer.Close()

	client := api.NewClient(mockServer.URL, "test-api-key")
	ctx := context.Background()

	t.Run("Events", func(t *testing.T) {
		t.Run("List", func(t *testing.T) {
			testListEvents(ctx, client)(t)
		})
		t.Run("Create", func(t *testing.T) {
			testCreateEvent(ctx, client)(t)
		})
		t.Run("Update", func(t *testing.T) {
			testUpdateEvent(ctx, client)(t)
		})
		t.Run("GetByName", func(t *testing.T) {
			testGetEventByName(ctx, client)(t)
		})
	})

	t.Run("AccessKeys", func(t *testing.T) {
		t.Run("List", func(t *testing.T) {
			testListAccessKeys(ctx, client)(t)
		})
		t.Run("Create", func(t *testing.T) {
			testCreateAccessKey(ctx, client)(t)
		})
		t.Run("GetPermissions", func(t *testing.T) {
			testGetAccessKeyPermissions(ctx, client)(t)
		})
		t.Run("SetPermissions", func(t *testing.T) {
			testSetAccessKeyPermissions(ctx, client)(t)
		})
	})
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
		event, err := client.GetEventByName(ctx, "test-event")
		require.NoError(t, err)
		require.NotNil(t, event)
		assert.Equal(t, "test-event", event.Name)
	}
}

func testCreateEvent(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		event := &domain.Event{
			Name: "new-event",
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
		permissions := &domain.Permissions{
			Send:    []string{"event1", "event2"},
			Receive: []string{"event3"},
		}
		created, err := client.CreateAccessKey(ctx, permissions)
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.AccessKey)
	}
}

func testGetAccessKeyPermissions(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		permissions, err := client.GetAccessKeyPermissions(ctx, "test-key")
		require.NoError(t, err)
		require.NotNil(t, permissions)
		assert.NotNil(t, permissions.Permissions)
	}
}

func testSetAccessKeyPermissions(ctx context.Context, client *api.Client) func(*testing.T) {
	return func(t *testing.T) {
		permissions := &domain.Permissions{
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
			{ID: 1, Name: "event1", Payload: map[string]string{"key": "value1"}},
			{ID: 2, Name: "event2", Payload: map[string]string{"key": "value2"}},
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

func mockGetEventByName(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "GET", "test-api-key")
	event := &domain.Event{
		ID:      1,
		Name:    "test-event",
		Payload: map[string]string{"key": "value"},
	}
	sendJSONResponse(w, event)
}

func mockListAccessKeys(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "GET", "test-api-key")
	sendJSONResponse(w, domain.AccessKeyList{
		ResultsLength: 2,
		Results: []*domain.AccessKeyPermissions{
			{
				Key: "key1",
				Permissions: &domain.Permissions{
					Send:    []string{"event1"},
					Receive: []string{"event2"},
				},
			},
			{
				Key: "key2",
				Permissions: &domain.Permissions{
					Send:    []string{"event3"},
					Receive: []string{"event4"},
				},
			},
		},
	})
}

func mockCreateAccessKey(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "POST", "test-api-key")
	sendJSONResponse(w, map[string]string{"accessKey": "new-access-key-123"})
}

func mockGetAccessKeyPermissions(t *testing.T, w http.ResponseWriter, r *http.Request) {
	verifyRequest(t, r, "GET", "test-api-key")
	sendJSONResponse(w, &domain.AccessKeyPermissions{
		Key: "test-key",
		Permissions: &domain.Permissions{
			Send:    []string{"event1"},
			Receive: []string{"event2"},
		},
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
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
