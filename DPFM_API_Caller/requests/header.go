package requests

type Header struct {
	InvoiceDocument     	int     `json:"InvoiceDocument"`
	HeaderBillingConfStatus *string `json:"HeaderBillingConfStatus"`
	IsCancelled         	*bool   `json:"IsCancelled"`
}
