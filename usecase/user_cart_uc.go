package usecase

import (
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/pkg/str"
	"retel-backend/usecase/viewmodel"
	"strings"
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
