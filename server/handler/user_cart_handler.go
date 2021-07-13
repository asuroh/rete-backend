package handler

import (
	"net/http"
	"retel-backend/server/request"
	"retel-backend/usecase"
	"strconv"

	"github.com/go-chi/chi"
	validator "gopkg.in/go-playground/validator.v9"
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

// CheckoutHandler ...
func (h *UserCartHandler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterfaceWithNil(r.Context(), "user")
	userID := user["id"].(string)

	req := request.UserCartRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	userCartUC := usecase.UserCartUC{ContractUC: h.ContractUC}
	res, err := userCartUC.ToCart(userID, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// UpdateHandler ...
func (h *UserCartHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	req := request.UserCartUpdateRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	userCartUC := usecase.UserCartUC{ContractUC: h.ContractUC}
	res, err := userCartUC.Update(id, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// DeleteHandler ...
func (h *UserCartHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	userCartUC := usecase.UserCartUC{ContractUC: h.ContractUC}
	err := userCartUC.Delete(id)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, "success", nil)
	return
}
