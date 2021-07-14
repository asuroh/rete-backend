package usecase

import (
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/pkg/str"
	"retel-backend/usecase/viewmodel"
	"strings"
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
