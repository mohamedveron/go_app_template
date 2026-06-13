package serverV1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api_gen "github.com/mohamedveron/go_app_template/internal/transport/http/api_server_gen/v1"
	"github.com/mohamedveron/go_app_template/internal/users"
	"github.com/mohamedveron/go_app_template/internal/users/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakePersistence is an in-memory implementation of persistence.UsersPersistence.
type fakePersistence struct {
	users  []*domain.User
	forceErr error
}

func (f *fakePersistence) Create(_ context.Context, u *domain.User) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	for _, existing := range f.users {
		if existing.Email == u.Email {
			return errors.New("user with that email already exists")
		}
	}
	f.users = append(f.users, u)
	return nil
}

func (f *fakePersistence) ReadByEmail(_ context.Context, email string) (*domain.User, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	for _, u := range f.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("email not found")
}

func (f *fakePersistence) List(_ context.Context, cursor domain.Cursor) (*domain.UserPage, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}

	limit := cursor.Limit
	if limit == 0 {
		limit = 20
	}

	var filtered []*domain.User
	for _, u := range f.users {
		if cursor.After == nil || (u.CreatedAt != nil && u.CreatedAt.After(*cursor.After)) {
			filtered = append(filtered, u)
		}
	}

	page := &domain.UserPage{}
	if uint(len(filtered)) > limit {
		page.Users = filtered[:limit]
		page.NextCursor = filtered[limit].CreatedAt
	} else {
		page.Users = filtered
	}
	return page, nil
}

func newTestHTTP(p *fakePersistence) *HTTP {
	svc, err := users.NewService(p)
	if err != nil {
		panic(err)
	}
	return &HTTP{users: svc}
}

func ptr[T any](v T) *T { return &v }

// makeTime returns a deterministic time offset by n seconds from a base.
func makeTime(n int) time.Time {
	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	return base.Add(time.Duration(n) * time.Second)
}

func seedUsers(n int) []*domain.User {
	out := make([]*domain.User, n)
	for i := range n {
		t := makeTime(i)
		out[i] = &domain.User{
			FirstName: "User",
			LastName:  "Test",
			Email:     "user" + string(rune('a'+i)) + "@example.com",
			CreatedAt: &t,
			UpdatedAt: &t,
		}
	}
	return out
}

// ── AddUser ──────────────────────────────────────────────────────────────────

func TestAddUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       any
		forceErr   error
		wantStatus int
		wantEmail  string
	}{
		{
			name:       "creates user successfully",
			body:       api_gen.NewUser{FirstName: "Jane", Email: "jane@example.com"},
			wantStatus: http.StatusCreated,
			wantEmail:  "jane@example.com",
		},
		{
			name:       "duplicate email",
			body:       api_gen.NewUser{FirstName: "Jane", Email: "jane@example.com"},
			forceErr:   errors.New("user with that email already exists"),
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "invalid JSON body",
			body:       "not-json{{",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email format",
			body:       api_gen.NewUser{FirstName: "Jane", Email: "not-an-email"},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &fakePersistence{forceErr: tt.forceErr}
			h := newTestHTTP(p)

			bodyBytes, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.AddUser(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantEmail != "" {
				var got api_gen.User
				require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
				require.NotNil(t, got.Email)
				assert.Equal(t, tt.wantEmail, *got.Email)
			}
		})
	}
}

// ── ListUsers ─────────────────────────────────────────────────────────────────

func TestListUsers(t *testing.T) {
	t.Parallel()

	seeded := seedUsers(5)

	tests := []struct {
		name           string
		storedUsers    []*domain.User
		params         api_gen.ListUsersParams
		forceErr       error
		wantStatus     int
		wantCount      int
		wantNextCursor bool
	}{
		{
			name:        "returns all users when store has fewer than default limit",
			storedUsers: seeded,
			params:      api_gen.ListUsersParams{},
			wantStatus:  http.StatusOK,
			wantCount:   5,
		},
		{
			name:        "respects limit param",
			storedUsers: seeded,
			params:      api_gen.ListUsersParams{Limit: ptr[int32](2)},
			wantStatus:  http.StatusOK,
			wantCount:   2,
			wantNextCursor: true,
		},
		{
			name:        "cursor skips already-seen rows",
			storedUsers: seeded,
			// after the 3rd user's createdAt → should return users 4 and 5
			params:     api_gen.ListUsersParams{After: seeded[2].CreatedAt},
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:        "empty store returns empty list",
			storedUsers: nil,
			params:      api_gen.ListUsersParams{},
			wantStatus:  http.StatusOK,
			wantCount:   0,
		},
		{
			name:        "persistence error returns 500",
			storedUsers: seeded,
			forceErr:    errors.New("db unavailable"),
			wantStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := &fakePersistence{users: tt.storedUsers, forceErr: tt.forceErr}
			h := newTestHTTP(p)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
			w := httptest.NewRecorder()

			h.ListUsers(w, req, tt.params)

			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				var got api_gen.UserPage
				require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
				assert.Len(t, got.Users, tt.wantCount)
				if tt.wantNextCursor {
					assert.NotNil(t, got.NextCursor)
				} else {
					assert.Nil(t, got.NextCursor)
				}
			}
		})
	}
}
