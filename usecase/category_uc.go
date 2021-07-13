package usecase

import (
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/pkg/str"
	"retel-backend/usecase/viewmodel"
	"strings"
)

// CategoryUC ...
type CategoryUC struct {
	*ContractUC
}

// BuildBody ...
func (uc CategoryUC) BuildBody(data *model.CategoryEntity, res *viewmodel.CategoryVM) {
	res.ID = data.ID
	res.Name = data.Name
	res.Description = data.Description
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt.String
}

// FindAll ...
func (uc CategoryUC) FindAll(search string, page, limit int, by, sort string) (res []viewmodel.CategoryVM, pagination viewmodel.PaginationVM, err error) {
	ctx := "CategoryUC.FindAll"

	if !str.Contains(model.CategoryBy, by) {
		by = model.DefaultCategoryBy
	}
	if !str.Contains(SortWhitelist, strings.ToLower(sort)) {
		sort = DescSort
	}

	limit = uc.LimitMax(limit)
	limit, offset := uc.PaginationPageOffset(page, limit)

	m := model.NewCategoryModel(uc.DB)
	data, count, err := m.FindAll(search, offset, limit, by, sort)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, pagination, err
	}
	pagination = PaginationRes(page, count, limit)

	for _, r := range data {
		temp := viewmodel.CategoryVM{}
		uc.BuildBody(&r, &temp)
		res = append(res, temp)
	}

	return res, pagination, err
}
