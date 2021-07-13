package usecase

import (
	"errors"
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/pkg/str"
	"retel-backend/server/request"
	"retel-backend/usecase/viewmodel"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// UserCartUC ...
type UserCartUC struct {
	*ContractUC
}

// BuildBody ...
func (uc UserCartUC) BuildBody(data *model.UserCartEntity, res *viewmodel.UserCartVM) {
	res.ID = data.ID
	res.UserID = data.UserID
	res.ProductID = data.ProductID
	res.Qty = data.Qty
	res.Price = data.Price
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt.String
}

// FindAll ...
func (uc UserCartUC) FindAll(userID string, page, limit int, by, sort string) (res []viewmodel.UserCartVM, pagination viewmodel.PaginationVM, err error) {
	ctx := "UserCartUC.FindAll"

	if !str.Contains(model.UserCartBy, by) {
		by = model.DefaultUserCartBy
	}
	if !str.Contains(SortWhitelist, strings.ToLower(sort)) {
		sort = DescSort
	}

	limit = uc.LimitMax(limit)
	limit, offset := uc.PaginationPageOffset(page, limit)

	m := model.NewUserCartModel(uc.DB)
	data, count, err := m.FindAll(userID, offset, limit, by, sort)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, pagination, err
	}
	pagination = PaginationRes(page, count, limit)

	for _, r := range data {
		temp := viewmodel.UserCartVM{}
		uc.BuildBody(&r, &temp)
		res = append(res, temp)
	}

	return res, pagination, err
}

// CheckDetails ...
func (uc UserCartUC) CheckDetails(data *request.UserCartRequest, product *viewmodel.ProductVM) (err error) {
	ctx := "UserCartUC.CheckDetails"

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

// ToCart ...
func (uc UserCartUC) ToCart(userID string, data *request.UserCartRequest) (res viewmodel.UserCartVM, err error) {
	ctx := "UserCartUC.ToCart"
	var product viewmodel.ProductVM

	err = uc.CheckDetails(data, &product)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.UserCartVM{
		ID:        uuid.NewV4().String(),
		ProductID: data.ProductID,
		UserID:    userID,
		Qty:       data.Qty,
		Price:     product.Price,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	m := model.NewUserCartModel(uc.DB)
	err = m.Store(res.ID, res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// Update ...
func (uc UserCartUC) Update(id string, data *request.UserCartUpdateRequest) (res viewmodel.UserCartVM, err error) {
	ctx := "UserCartUC.Update"
	now := time.Now().UTC()
	res = viewmodel.UserCartVM{
		ID:        id,
		Qty:       data.Qty,
		UpdatedAt: now.Format(time.RFC3339),
	}
	m := model.NewUserCartModel(uc.DB)
	err = m.Update(id, res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// Delete ...
func (uc UserCartUC) Delete(id string) (err error) {
	ctx := "UserCartUC.Delete"

	now := time.Now().UTC()
	m := model.NewUserCartModel(uc.DB)
	err = m.Destroy(id, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return err
	}

	return err
}
