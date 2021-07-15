package usecase

import (
	"errors"
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/server/request"
	"retel-backend/usecase/viewmodel"
	"time"

	uuid "github.com/satori/go.uuid"
)

// TransactionDetailUC ...
type TransactionDetailUC struct {
	*ContractUC
}

// CheckDetails ...
func (uc TransactionDetailUC) CheckDetails(data *request.TransactionDetailRequest, product *viewmodel.ProductVM) (err error) {
	ctx := "TransactionDetailUC.CheckDetails"

	productUC := ProductUC{ContractUC: uc.ContractUC}
	p, err := productUC.FindByID(data.ProductID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "FindByID", uc.ReqID)
		return err
	}

	if p.ID == "" {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "check_product", uc.ReqID)
		return errors.New("invalid_product")
	}

	if data.Qty > int64(p.Qty) {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "check_qty", uc.ReqID)
		return errors.New("invalid_qty")
	}

	product.Price = p.Price

	return err
}

// Create ...
func (uc TransactionDetailUC) Create(data *request.TransactionDetailRequest) (res viewmodel.TransactionDetailVM, err error) {
	ctx := "TransactionDetailUC.Create"
	var product viewmodel.ProductVM

	err = uc.CheckDetails(data, &product)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.TransactionDetailVM{
		ID:            uuid.NewV4().String(),
		TransactionID: data.TransactionID,
		ProductName:   product.Name,
		ProductID:     data.ProductID,
		Price:         product.Price,
		Qty:           data.Qty,
		CheckIn:       data.CheckIn,
		CheckOut:      data.CheckOut,
		Total:         product.Price * float64(data.Qty),
		CreatedAt:     now.Format(time.RFC3339),
		UpdatedAt:     now.Format(time.RFC3339),
	}
	m := model.NewTransactionDetailModel(uc.DB)
	err = m.Store(res.ID, res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query_detail", uc.ReqID)
		return res, err
	}

	return res, err
}
