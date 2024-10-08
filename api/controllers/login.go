package controllers

import (
	"myapp/internal/users"
	"myapp/pkg/io"
	"net/http"
)

func (h *HandlerLayer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req users.LoginRequest
	err := io.ReadJSON(r, &req)
	if err != nil {
		io.SendError(w, "Bad Request", http.StatusBadRequest)
		return
	}
	user, err := h.LogicLayer.VerifyUserCredentials(req.Username, req.Password)
	if err != nil {
		io.SendError(w, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := h.LogicLayer.CreateSession(user.Username)
	if err != nil {
		io.SendError(w, "Failed to create (rewrite) session", http.StatusUnauthorized)
		h.Logger.Errorf("Error in create session: %v", err)
		return
	}
	h.Logger.Infof("Logined user: %s", user.Username)
	io.WriteJSON(w, http.StatusOK, users.LoginResponse{Token: token})
}
