package handler

import (
	"net/http"
	"retel-backend/usecase"
	"strconv"
)

// CategoryHandler ...
type CategoryHandler struct {
	Handler
}

// GetAllHandler ...
func (h *CategoryHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
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
	search := r.URL.Query().Get("search")
	by := r.URL.Query().Get("by")
	sort := r.URL.Query().Get("sort")

	categoryUC := usecase.CategoryUC{ContractUC: h.ContractUC}
	res, p, err := categoryUC.FindAll(search, page, limit, by, sort)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, p)
	return
}
