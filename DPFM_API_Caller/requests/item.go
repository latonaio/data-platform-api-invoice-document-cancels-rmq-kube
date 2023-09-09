package requests

type Item struct {
	InvoiceDocument     	int     `json:"InvoiceDocument"`
	InvoiceDocumentItem 	int     `json:"InvoiceDocumentItem"`
	ItemBillingConfStatus   *string `json:"ItemBillingConfStatus"`
	IsCancelled         	*bool   `json:"IsCancelled"`
}
