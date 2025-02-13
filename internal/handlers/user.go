package handlers

import "net/http"

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) error {
	return render(w, http.StatusOK, "create-user.html", nil)
}
