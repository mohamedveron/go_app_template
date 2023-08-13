package http

import "net/http"

// AddUser implements ServerInterface.
func (*HTTP) AddUser(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// FindUserByID implements ServerInterface.
func (*HTTP) FindUserByID(w http.ResponseWriter, r *http.Request, id int64) {
	panic("unimplemented")
}
