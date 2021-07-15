package usecase

import (
	"errors"
	"fmt"
	"retel-backend/model"
	"retel-backend/pkg/logruslogger"
	"retel-backend/server/request"
	"retel-backend/usecase/viewmodel"
)

// XenditUC ...
type XenditUC struct {
	*ContractUC
}

// XenditInvoice ...
func (uc XenditUC) XenditInvoice(code, userID, transactionID string, total float64) (res viewmodel.InvoiceData, err error) {
	ctx := "XenditUC.XenditInvoice"

	userUc := UserUC{ContractUC: uc.ContractUC}
	user, err := userUc.FindByID(userID, false)

	bodyXendit := request.XenditInvoiceRequest{
		ExternalID:         code,
		PayerEmail:         user.Email,
		Description:        "Pembayaran Pembelanjaan Di Belanja Parts",
		ShouldSendEmail:    true,
		Amount:             total,
		SuccessRedirectUrl: uc.ContractUC.EnvConfig["SUCCESS_REDIRECT_URL"],
		FailureRedirectUrl: uc.ContractUC.EnvConfig["FAILURE_REDIRECT_URL"],
	}

	res, err = uc.Xendit.CreateInvoice(bodyXendit)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "create_xendit", uc.ReqID)
		return res, err
	}
	return res, err
}

// XenditInvoiceCallback ...
func (uc XenditUC) XenditInvoiceCallback(data request.XenditInvoiceCallbackRequest) (res string, err error) {
	ctx := "XenditUC.XenditInvoice"

	transactionUc := TransactionUC{ContractUC: uc.ContractUC}
	transactionData, err := transactionUc.FindByInvoiceID(data.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_transaction", uc.ReqID)
		return res, errors.New("Transaksi not Found")
	}

	Status := ""
	PaymentMethod := ""
	transaction, err := uc.Xendit.GetInvoice(transactionData.InvoceID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "get_invoice_xendit", uc.ReqID)
		return res, err
	}

	if transaction.Status == "PAID" || transaction.Status == "SETTLED" {
		Status = model.TransactionPaid
		PaymentMethod = data.PaymentChannel
	} else if transaction.Status == "EXPIRED" {
		Status = model.TransactionCanceled
	} else {
		Status = model.TransactionCanceled
	}
	fmt.Println(PaymentMethod)
	_, err = transactionUc.Update(transactionData.ID, Status, transactionData.Status)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "update_status", uc.ReqID)
		return res, err
	}

	return res, err
}
