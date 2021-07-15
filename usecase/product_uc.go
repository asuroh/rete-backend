package usecase

import (
	"errors"
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/pkg/str"
	"retel-backend/usecase/viewmodel"
	"strings"
	"time"
)

// ProductUC ...
type ProductUC struct {
	*ContractUC
}

// BuildBody ...
func (uc ProductUC) BuildBody(data *model.ProductEntity, res *viewmodel.ProductVM) {

	res.ID = data.ID
	res.Name = data.Name
	res.CategoryID = data.CategoryID
	res.CategoryName = data.CategoryName
	res.Description = data.Description
	res.Price = data.Price
	res.Qty = data.Qty
	if data.ImagePath.String != "" {
		res.ImagePath = uc.EnvConfig["APP_IMAGE_URL"] + uc.EnvConfig["FILE_PATH"] + data.ImagePath.String
	}
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt.String
}

// FindAll ...
func (uc ProductUC) FindAll(search, categoryID string, page, limit int, by, sort string) (res []viewmodel.ProductVM, pagination viewmodel.PaginationVM, err error) {
	ctx := "ProductUC.FindAll"

	if !str.Contains(model.UserBy, by) {
		by = model.DefaultUserBy
	}
	if !str.Contains(SortWhitelist, strings.ToLower(sort)) {
		sort = DescSort
	}

	limit = uc.LimitMax(limit)
	limit, offset := uc.PaginationPageOffset(page, limit)

	m := model.NewProductModel(uc.DB)
	data, count, err := m.FindAll(search, categoryID, offset, limit, by, sort)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, pagination, err
	}
	pagination = PaginationRes(page, count, limit)

	for _, r := range data {
		temp := viewmodel.ProductVM{}
		uc.BuildBody(&r, &temp)
		res = append(res, temp)
	}

	return res, pagination, err
}

// FindByID ...
func (uc ProductUC) FindByID(id string) (res viewmodel.ProductVM, err error) {
	ctx := "ProductUC.FindByID"

	m := model.NewProductModel(uc.DB)
	data, err := m.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}
	uc.BuildBody(&data, &res)

	return res, err
}

// CheckDetails ...
func (uc ProductUC) CheckDetails(oldData *viewmodel.ProductVM, product *viewmodel.ProductVM) (err error) {
	ctx := "ProductUC.CheckDetails"

	p, err := uc.FindByID(oldData.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "FindByID", uc.ReqID)
		return err
	}

	if p.ID == "" {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "check_product", uc.ReqID)
		return errors.New("invalid_product")
	}

	if int64(oldData.Qty) > int64(p.Qty) || p.Qty == 0 {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "check_qty", uc.ReqID)
		return errors.New("invalid_qty")
	}
	product.Qty = p.Qty
	return err
}

// UpdateStock ...
func (uc ProductUC) UpdateStock(id string, typeData string, qty int64) (err error) {
	ctx := "ProductUC.UpdateStock"
	var product viewmodel.ProductVM
	var qtyReal int64

	err = uc.CheckDetails(&viewmodel.ProductVM{
		ID:  id,
		Qty: int(qty),
	}, &product)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_details", uc.ReqID)
		return err
	}

	if typeData == model.TypeDataMinus {
		qtyReal = int64(product.Qty) - qty
	} else if typeData == model.TypeDataPlus {
		qtyReal = int64(product.Qty) + qty
	} else {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "type_data", uc.ReqID)
		return errors.New("invalid_type")
	}

	now := time.Now().UTC()
	m := model.NewProductModel(uc.DB)
	err = m.UpdateStock(id, qtyReal, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return err
	}
	return err
}
