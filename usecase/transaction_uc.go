package usecase

import (
	"errors"
	"retel-backend/helper"
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/pkg/str"
	"retel-backend/server/request"
	"retel-backend/usecase/viewmodel"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// TransactionUC ...
type TransactionUC struct {
	*ContractUC
}

// BuildBody ...
func (uc TransactionUC) BuildBody(data *model.TransactionEntity, res *viewmodel.TransactionVM) {
	res.ID = data.ID
	res.UserID = data.UserID
	res.Code = data.Code
	res.Total = data.Total
	res.Note = data.Note.String
	res.UrlPayment = data.UrlPayment.String
	res.InvoceID = data.InvoceID.String
	res.Status = data.Status
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt.String
}

// FindAll ...
func (uc TransactionUC) FindAll(userID string, page, limit int, by, sort string) (res []viewmodel.TransactionVM, pagination viewmodel.PaginationVM, err error) {
	ctx := "TransactionUC.FindAll"

	if !str.Contains(model.TransactionBy, by) {
		by = model.DefaultTransactionBy
	}
	if !str.Contains(SortWhitelist, strings.ToLower(sort)) {
		sort = DescSort
	}

	limit = uc.LimitMax(limit)
	limit, offset := uc.PaginationPageOffset(page, limit)

	m := model.NewTransactionModel(uc.DB)
	data, count, err := m.FindAll(userID, offset, limit, by, sort)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, pagination, err
	}
	pagination = PaginationRes(page, count, limit)

	for _, r := range data {
		temp := viewmodel.TransactionVM{}
		uc.BuildBody(&r, &temp)
		res = append(res, temp)
	}

	return res, pagination, err
}

// FindByID ...
func (uc TransactionUC) FindByID(id string) (res viewmodel.TransactionVM, err error) {
	ctx := "TransactionUC.FindByID"

	m := model.NewTransactionModel(uc.DB)
	data, err := m.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}
	uc.BuildBody(&data, &res)

	return res, err
}

// FindByInvoiceID ...
func (uc TransactionUC) FindByInvoiceID(invoceID string) (res viewmodel.TransactionVM, err error) {
	ctx := "TransactionUC.FindByInvoiceID"

	m := model.NewTransactionModel(uc.DB)
	data, err := m.FindByInvoiceID(invoceID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}
	uc.BuildBody(&data, &res)

	return res, err
}

// Create ...
func (uc TransactionUC) Create(data *request.TransactionRequest) (res viewmodel.TransactionVM, err error) {
	ctx := "TransactionUC.Create"
	m := model.NewTransactionModel(uc.DB)

	count, err := m.CountNumberInvoice()
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query_count", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.TransactionVM{
		ID:        uuid.NewV4().String(),
		UserID:    data.UserID,
		Note:      data.Note,
		Code:      helper.GenerateInvoice(count),
		Status:    model.TransactionWaitingPayment,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	err = m.Store(res.ID, res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	transactionDetailUC := TransactionDetailUC{ContractUC: uc.ContractUC}
	for _, row := range data.TransactionDetail {
		row.TransactionID = res.ID
		detail, err := transactionDetailUC.Create(&row)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "store_detail", uc.ReqID)
			return res, err
		}
		res.Total = res.Total + detail.Total
		res.TransactionDetail = append(res.TransactionDetail, detail)
	}

	xenditUC := XenditUC{ContractUC: uc.ContractUC}
	xendit, err := xenditUC.XenditInvoice(res.Code, res.UserID, res.ID, res.Total)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "xendit", uc.ReqID)
		return res, err
	}

	res.InvoceID = xendit.ID
	res.UrlPayment = xendit.InvoiceUrl
	err = m.UpdateXendit(res.ID, viewmodel.TransactionXenditVM{
		Total:      res.Total,
		UrlPayment: res.UrlPayment,
		InvoceID:   res.InvoceID,
	}, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "update_url", uc.ReqID)
		return res, err
	}

	return res, err
}

// Update ...
func (uc TransactionUC) Update(id, status, oldStatus string) (res viewmodel.TransactionVM, err error) {
	ctx := "UserUC.Update"

	if oldStatus != model.TransactionWaitingPayment {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "check_status_old", uc.ReqID)
		return res, errors.New("invalid_status")
	}

	now := time.Now().UTC()
	res = viewmodel.TransactionVM{
		ID:        id,
		Status:    status,
		UpdatedAt: now.Format(time.RFC3339),
	}
	m := model.NewTransactionModel(uc.DB)
	err = m.UpdateStatus(id, status, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}
