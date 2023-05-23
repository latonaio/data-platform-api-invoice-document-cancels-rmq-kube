package requests

type ScheduleLine struct {
	InvoiceDocumentID  int   `json:"InvoiceDocumentID"`
	InvoiceDocument    int   `json:"InvoiceDocument"`
	ScheduleLine       int   `json:"ScheduleLine"`
	IsCancelled        *bool `json:"IsCancelled"`
	IsMarkedForDeleted *bool `json:"IsMarkedForDeleted"`
}
