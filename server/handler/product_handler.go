package handler

import (
	"net/http"
	"retel-backend/usecase"
	"strconv"

	"github.com/go-chi/chi"
)

// ProductHandler ...
type ProductHandler struct {
	Handler
}

// GetAllHandler ...
func (h *ProductHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
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
	categoryID := r.URL.Query().Get("category_id")
	search := r.URL.Query().Get("search")
	by := r.URL.Query().Get("by")
	sort := r.URL.Query().Get("sort")

	productUC := usecase.ProductUC{ContractUC: h.ContractUC}
	res, p, err := productUC.FindAll(search, categoryID, page, limit, by, sort)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, p)
	return
}

// GetByIDHandler ...
func (h *ProductHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	productUC := usecase.ProductUC{ContractUC: h.ContractUC}
	res, err := productUC.FindByID(id)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
