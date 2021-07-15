package xendit

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"retel-backend/server/request"
	"retel-backend/usecase/viewmodel"
)

// Credential ...
type Credential struct {
	Key string
}

var (
	postURL = "https://api.xendit.co"
)

// CreateInvoice ...
func (m Credential) CreateInvoice(data request.XenditInvoiceRequest) (res viewmodel.InvoiceData, err error) {
	endPoint := postURL + `/v2/invoices`

	invoiceDataInBytes := new(bytes.Buffer)
	json.NewEncoder(invoiceDataInBytes).Encode(data)

	req, err := http.NewRequest("POST", endPoint, invoiceDataInBytes)
	req.SetBasicAuth(m.Key, "")
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return res, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		dataError := viewmodel.InvoiceError{}
		bodyError, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(bodyError, &dataError)

		return res, errors.New(dataError.Message)
	}

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	return res, err
}

// GetInvoice ...
func (m Credential) GetInvoice(invoiceID string) (res viewmodel.InvoiceData, err error) {
	endPoint := postURL + `/v2/invoices/` + invoiceID

	req, err := http.NewRequest("GET", endPoint, bytes.NewBuffer(nil))
	req.SetBasicAuth(m.Key, "")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return res, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		dataError := viewmodel.InvoiceError{}
		bodyError, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(bodyError, &dataError)

		return res, errors.New(dataError.Message)
	}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	return res, err
}
