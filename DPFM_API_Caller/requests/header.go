package requests

type Header struct {
	InvoiceDocumentID   int     `json:"InvoiceDocumentID"`
	HeaderInvoiceStatus *string `json:"HeaderInvoiceStatus"`
	IsCancelled         *bool   `json:"IsCancelled"`
	HeaderIsDeleted     *bool   `json:"HeaderIsDeleted"`
}
