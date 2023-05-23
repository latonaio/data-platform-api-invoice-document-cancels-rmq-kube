package requests

type Item struct {
	InvoiceDocumentID int     `json:"InvoiceDocumentID"`
	InvoiceDocument   int     `json:"InvoiceDocument"`
	ItemInvoiceStatus *string `json:"ItemInvoiceStatus"`
	IsCancelled       *bool   `json:"IsCancelled"`
	ItemIsDeleted     *bool   `json:"ItemIsDeleted"`
}
