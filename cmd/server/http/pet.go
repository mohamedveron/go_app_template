package http

import "net/http"

// AddPet implements ServerInterface.
func (*HTTP) AddPet(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

// DeletePet implements ServerInterface.
func (*HTTP) DeletePet(w http.ResponseWriter, r *http.Request, id int64) {
	panic("unimplemented")
}

// FindPetByID implements ServerInterface.
func (*HTTP) FindPetByID(w http.ResponseWriter, r *http.Request, id int64) {
	panic("unimplemented")
}
