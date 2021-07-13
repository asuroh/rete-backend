package handler

import (
	"net/http"
	"retel-backend/usecase"
	"strconv"
)

// UserCartHandler ...
type UserCartHandler struct {
	Handler
}

// GetAllHandler ...
func (h *UserCartHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterfaceWithNil(r.Context(), "user")
	userID := user["id"].(string)

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		SendBadRequest(w, "Invalid page value")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		SendBadRequest(w, "Invalid limit value")
		return
	}
	by := r.URL.Query().Get("by")
	sort := r.URL.Query().Get("sort")

	userCartUC := usecase.UserCartUC{ContractUC: h.ContractUC}
	res, p, err := userCartUC.FindAll(userID, page, limit, by, sort)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, p)
	return
}
