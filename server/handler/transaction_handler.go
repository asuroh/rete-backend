package handler

import (
	"net/http"
	"retel-backend/server/request"
	"retel-backend/usecase"
	"strconv"

	"github.com/go-chi/chi"
	validator "gopkg.in/go-playground/validator.v9"
)

// TransactionHandler ...
type TransactionHandler struct {
	Handler
}

// GetAllByTokenHandler ...
func (h *TransactionHandler) GetAllByTokenHandler(w http.ResponseWriter, r *http.Request) {
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

	transactionUC := usecase.TransactionUC{ContractUC: h.ContractUC}
	res, p, err := transactionUC.FindAll(userID, page, limit, by, sort)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, p)
	return
}

// GetByIDHandler ...
func (h *TransactionHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	transactionUC := usecase.TransactionUC{ContractUC: h.ContractUC}
	res, err := transactionUC.FindByID(id)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// CheckoutHandler ...
func (h *TransactionHandler) CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterfaceWithNil(r.Context(), "user")
	userID := user["id"].(string)

	req := request.TransactionRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}
	req.UserID = userID
	transactionUC := usecase.TransactionUC{ContractUC: h.ContractUC}
	res, err := transactionUC.Create(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

//XenditInvoiceCallbackHandler ...
func (h *TransactionHandler) XenditInvoiceCallbackHandler(w http.ResponseWriter, r *http.Request) {
	req := request.XenditInvoiceCallbackRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	xenditUc := usecase.XenditUC{ContractUC: h.ContractUC}
	res, err := xenditUc.XenditInvoiceCallback(req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
