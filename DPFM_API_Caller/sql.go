package dpfm_api_caller

import (
	dpfm_api_input_reader "data-platform-api-invoice-document-cancels-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-invoice-document-cancels-rmq-kube/DPFM_API_Output_Formatter"

	"fmt"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
)

func (c *DPFMAPICaller) HeaderRead(
	input *dpfm_api_input_reader.SDC,
	log *logger.Logger,
) *dpfm_api_output_formatter.Header {
	where := fmt.Sprintf("WHERE header.InvoiceDocument = %d ", input.InvoiceDocument.InvoiceDocument)
	where = fmt.Sprintf("%s \n AND ( header.BillToParty = %d OR header.BillFromParty = %d ) ", where, input.BusinessPartner, input.BusinessPartner)
	// where = fmt.Sprintf("%s \n AND ( header.IsCancelled ) = ( false ) ", where)
	rows, err := c.db.Query(
		`SELECT 
			header.InvoiceDocument
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_invoice_document_header_data as header ` + where + ` ;`)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToHeader(rows)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}

	return data
}

func (c *DPFMAPICaller) ItemsRead(
	input *dpfm_api_input_reader.SDC,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.Item {
	where := fmt.Sprintf("WHERE item.InvoiceDocument IS NOT NULL\nAND header.InvoiceDocument = %d", input.InvoiceDocument.InvoiceDocument)
	where = fmt.Sprintf("%s \n AND ( header.BillToParty = %d OR header.BillFromParty = %d ) ", where, input.BusinessPartner, input.BusinessPartner)
	// where = fmt.Sprintf("%s \n AND ( item.IsCancelled ) = ( false ) ", where)
	rows, err := c.db.Query(
		`SELECT 
			item.InvoiceDocument, item.InvoiceDocumentItem
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_invoice_document_item_data as item
		INNER JOIN DataPlatformMastersAndTransactionsMysqlKube.data_platform_invoice_document_header_data as header
		ON header.InvoiceDocument = item.InvoiceDocument ` + where + ` ;`)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}
	defer rows.Close()

	data, err := dpfm_api_output_formatter.ConvertToItem(rows)
	if err != nil {
		log.Error("%+v", err)
		return nil
	}

	return data
}
