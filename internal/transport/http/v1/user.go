package serverV1

import (
	"encoding/json"
	"net/http"
	"time"

	api_gen "github.com/mohamedveron/go_app_template/internal/transport/http/api_server_gen/v1"
	"github.com/mohamedveron/go_app_template/internal/users/domain"
)

// ListUsers implements ServerInterface.
func (h *HTTP) ListUsers(w http.ResponseWriter, r *http.Request, params api_gen.ListUsersParams) {
	cursor := domain.Cursor{}
	if params.After != nil {
		cursor.After = params.After
	}
	if params.Limit != nil {
		cursor.Limit = uint(*params.Limit)
	}

	page, err := h.users.ListUsers(r.Context(), cursor)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := api_gen.UserPage{
		Users: make([]api_gen.User, 0, len(page.Users)),
	}
	for _, u := range page.Users {
		resp.Users = append(resp.Users, domainUserToAPI(u))
	}
	resp.NextCursor = page.NextCursor

	writeJSON(w, http.StatusOK, resp)
}

// AddUser implements ServerInterface.
func (h *HTTP) AddUser(w http.ResponseWriter, r *http.Request) {
	var body api_gen.AddUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	u := &domain.User{
		FirstName: body.FirstName,
		Email:     body.Email,
	}
	if body.LastName != nil {
		u.LastName = *body.LastName
	}
	if body.Mobile != nil {
		u.Mobile = *body.Mobile
	}

	created, err := h.users.CreateUser(r.Context(), u)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, domainUserToAPI(created))
}

// FindUserByID implements ServerInterface.
func (h *HTTP) FindUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	writeError(w, http.StatusNotImplemented, "not implemented")
}

func domainUserToAPI(u *domain.User) api_gen.User {
	au := api_gen.User{
		FirstName: &u.FirstName,
		LastName:  &u.LastName,
		Mobile:    &u.Mobile,
		Email:     &u.Email,
	}
	if u.CreatedAt != nil {
		t := time.Time(*u.CreatedAt)
		au.CreatedAt = &t
	}
	if u.UpdatedAt != nil {
		t := time.Time(*u.UpdatedAt)
		au.UpdatedAt = &t
	}
	return au
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, api_gen.Error{Code: int32(status), Message: msg})
}
