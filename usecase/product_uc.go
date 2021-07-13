package usecase

import (
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/pkg/str"
	"retel-backend/usecase/viewmodel"
	"strings"
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
