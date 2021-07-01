package handler

import (
	"io/ioutil"
	"net/http"
	"retel-backend/helper"
	"retel-backend/pkg/str"
	"retel-backend/server/request"
	"retel-backend/usecase"
	"strconv"

	"github.com/go-chi/chi"
	validator "gopkg.in/go-playground/validator.v9"
)

// UserHandler ...
type UserHandler struct {
	Handler
}

// LoginHandler ...
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := request.UserLoginRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUC.Login(req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// GetAllHandler ...
func (h *UserHandler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
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

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, p, err := userUC.FindAll(search, page, limit, by, sort)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, p)
	return
}

// GetByIDHandler ...
func (h *UserHandler) GetByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUC.FindByID(id, false)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// GetByTokenHandler ...
func (h *UserHandler) GetByTokenHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	id := user["id"].(string)

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUC.FindByID(id, false)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// RegisterHandler ...
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := request.UserRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUC.Register(&req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// UpdateByTokenHandler ...
func (h *UserHandler) UpdateByTokenHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	id := user["id"].(string)

	req := request.UserRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUC.Update(id, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// UpdateHandler ...
func (h *UserHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	req := request.UserRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUC.Update(id, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// UploadImageHandler ...
func (h *UserHandler) UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	req := request.UserUploadImageRequest{}
	user := requestIDFromContextInterfaceWithNil(r.Context(), "user")
	userID := user["id"].(string)

	maxUploadSize := str.StringToInt(h.Handler.EnvConfig["FILE_MAX_UPLOAD_SIZE"])
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxUploadSize))
	err := r.ParseMultipartForm(int64(maxUploadSize))
	if err != nil {
		SendBadRequest(w, helper.FileTooBig)
		return
	}

	file, header, err := r.FormFile("file")
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		SendBadRequest(w, helper.FileError)
		return
	}

	fileUploadUc := usecase.FileUploadUC{ContractUC: h.ContractUC}
	req.Path, err = fileUploadUc.Upload(header.Filename, req.Type, fileBytes)
	if err != nil {
		SendBadRequest(w, helper.FileError)
		return
	}

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUC.UpdateImage(userID, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// DeleteHandler ...
func (h *UserHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	userUC := usecase.UserUC{ContractUC: h.ContractUC}
	err := userUC.Delete(id)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, "success", nil)
	return
}
