package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-invoice-document-cancels-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-invoice-document-cancels-rmq-kube/DPFM_API_Output_Formatter"
	"data-platform-api-invoice-document-cancels-rmq-kube/config"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
	"golang.org/x/xerrors"
)

type DPFMAPICaller struct {
	ctx  context.Context
	conf *config.Conf
	rmq  *rabbitmq.RabbitmqClient
	db   *database.Mysql
}

func NewDPFMAPICaller(
	conf *config.Conf, rmq *rabbitmq.RabbitmqClient, db *database.Mysql,
) *DPFMAPICaller {
	return &DPFMAPICaller{
		ctx:  context.Background(),
		conf: conf,
		rmq:  rmq,
		db:   db,
	}
}

func (c *DPFMAPICaller) AsyncInvoiceDocumentCancels(
	accepter []string,
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) (interface{}, []error) {
	var response interface{}
	if input.APIType == "cancels" {
		response = c.cancelSqlProcess(input, output, accepter, log)
	}

	return response, nil
}

func (c *DPFMAPICaller) cancelSqlProcess(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	accepter []string,
	log *logger.Logger,
) *dpfm_api_output_formatter.Message {
	var headerData *dpfm_api_output_formatter.Header
	itemData := make([]dpfm_api_output_formatter.Item, 0)
	for _, a := range accepter {
		switch a {
		case "Header":
			h, i := c.headerCancel(input, output, log)
			if h == nil || i == nil {
				continue
			}
			headerData = h
			itemData = append(itemData, *i...)
		case "Item":
			i := c.itemCancel(input, output, log)
			if i == nil {
				continue
			}
			itemData = append(itemData, *i...)
		}
	}

	return &dpfm_api_output_formatter.Message{
		Header: headerData,
		Item:   &itemData,
	}
}

func (c *DPFMAPICaller) headerCancel(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) (*dpfm_api_output_formatter.Header, *[]dpfm_api_output_formatter.Item) {
	sessionID := input.RuntimeSessionID
	items := c.ItemsRead(input, log)
	for i := range *items {
		(*items)[i].IsCancelled = input.InvoiceDocument.IsCancelled
		res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": (*items)[i], "function": "InvoiceDocumentItem", "runtime_session_id": sessionID})
		if err != nil {
			err = xerrors.Errorf("rmq error: %w", err)
			log.Error("%+v", err)
			return nil, nil
		}
		res.Success()
		if !checkResult(res) {
			output.SQLUpdateResult = getBoolPtr(false)
			output.SQLUpdateError = "InvoiceDocument Item Data cannot cancel"
			return nil, nil
		}
	}

	header := c.HeaderRead(input, log)
	header.IsCancelled = input.InvoiceDocument.IsCancelled
	res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": header, "function": "InvoiceDocumentHeader", "runtime_session_id": sessionID})
	if err != nil {
		err = xerrors.Errorf("rmq error: %w", err)
		log.Error("%+v", err)
		return nil, nil
	}
	res.Success()
	if !checkResult(res) {
		output.SQLUpdateResult = getBoolPtr(false)
		output.SQLUpdateError = "Header Data cannot cancel"
		return nil, nil
	}
	return header, items
}

func (c *DPFMAPICaller) itemCancel(
	input *dpfm_api_input_reader.SDC,
	output *dpfm_api_output_formatter.SDC,
	log *logger.Logger,
) *[]dpfm_api_output_formatter.Item {
	sessionID := input.RuntimeSessionID

	items := make([]dpfm_api_output_formatter.Item, 0)
	for _, v := range input.InvoiceDocument.Item {
		data := dpfm_api_output_formatter.Item{
			InvoiceDocument:     input.InvoiceDocument.InvoiceDocument,
			InvoiceDocumentItem: v.InvoiceDocumentItem,
			ItemInvoiceStatus:   nil,
			IsCancelled:         v.IsCancelled,
			ItemIsDeleted:       nil,
		}
		res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": data, "function": "InvoiceDocumentItem", "runtime_session_id": sessionID})
		if err != nil {
			err = xerrors.Errorf("rmq error: %w", err)
			log.Error("%+v", err)
			return nil
		}
		res.Success()
		if !checkResult(res) {
			output.SQLUpdateResult = getBoolPtr(false)
			output.SQLUpdateError = "InvoiceDocument Item Data cannot cancel"
			return nil
		}
	}

	// itemがキャンセル取り消しされた場合、headerのキャンセルも取り消す
	if !*input.InvoiceDocument.Item[0].IsCancelled {
		header := c.HeaderRead(input, log)
		header.IsCancelled = input.InvoiceDocument.Item[0].IsCancelled
		res, err := c.rmq.SessionKeepRequest(nil, c.conf.RMQ.QueueToSQL()[0], map[string]interface{}{"message": header, "function": "InvoiceDocumentHeader", "runtime_session_id": sessionID})
		if err != nil {
			err = xerrors.Errorf("rmq error: %w", err)
			log.Error("%+v", err)
			return nil
		}
		res.Success()
		if !checkResult(res) {
			output.SQLUpdateResult = getBoolPtr(false)
			output.SQLUpdateError = "Header Data cannot cancel"
			return nil
		}
	}
	return &items
}

func checkResult(msg rabbitmq.RabbitmqMessage) bool {
	data := msg.Data()
	d, ok := data["result"]
	if !ok {
		return false
	}
	result, ok := d.(string)
	if !ok {
		return false
	}
	return result == "success"
}

func getBoolPtr(b bool) *bool {
	return &b
}
