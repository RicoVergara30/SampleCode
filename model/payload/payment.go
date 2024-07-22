package payload

import (
	"encoding/xml"
	"sample/model/bah"
)

type (
	PaymentRequest struct {
		RequestBody interface{}
	}
	PaymentStatusReport struct {
		PaymentStatusReport interface{}
	}
)

type (
	PaymentRequestISO20022 struct {
		XMLName xml.Name                       `xml:"Message"`
		Header  bah.HCRequestApplicationHeader `xml:"AppHdr"`
		Body    PaymentRequestWrapper
	}

	PaymentRequestWrapper struct {
		XMLName xml.Name           `xml:"PaymentRequest"`
		Body    PaymentRequestBody `xml:"CdtrPmtActvtnReq"`
	}
)

type (
	PaymentRequestBody struct {
		GroupHeader PaymentRequestGroupHeader `xml:"GrpHdr"`
		Information Information               `xml:"PmtInf"`
	}
)

type (
	PaymentRequestGroupHeader struct {
		MessageID            string `xml:"MsgId"`
		CreationDateTime     string `xml:"CreDtTm"`
		NumberOfTransactions string `xml:"NbOfTxs"`
		InitiatingBIC        string `xml:"InitgPty>Id>OrgId>AnyBIC"`
	}

	Information struct {
		ID                 string `xml:"PmtInfId"`
		Method             string `xml:"PmtMtd"`
		ExecutionDateTime  string `xml:"ReqdExctnDt>DtTm"`
		ExpirationDateTime string `xml:"XpryDt>DtTm"`
		DebtorName         string `xml:"Dbtr>Nm"`
		// DebtorAccount             DebtorAccount
		DebtorAgentBIC            string `xml:"DbtrAgt>FinInstnId>BICFI"`
		CreditTransferTransaction CreditTransferTransaction
	}

	CreditTransferTransaction struct {
		XMLName    xml.Name `xml:"CdtTrfTx"`
		EndToEndID string   `xml:"PmtId>EndToEndId"`
		// TypeInformation  PaymentTypeInformation
		Condition        Condition
		Amount           TransactionAmount `xml:"Amt>InstdAmt"`
		ChargeBearer     string            `xml:"ChrgBr"`
		CreditorAgentBIC string            `xml:"CdtrAgt>FinInstnId>BICFI"`
		CreditorName     string            `xml:"Cdtr>Nm"`
	}

	TransactionAmount struct {
		XMLName  xml.Name `xml:"InstdAmt"`
		Amount   string   `xml:",innerxml"`
		Currency string   `xml:"Ccy,attr"`
	}

	Condition struct {
		XMLName                   xml.Name `xml:"PmtCond"`
		AmountModAllowed          string   `xml:"AmtModAllwd"`
		EarlyPaymentAllowed       string   `xml:"EarlyPmtAllwd"`
		GeneratedPaymentRequested string   `xml:"GrntedPmtReqd"`
	}
)

type (
	TransactCredit struct {
		ReceivingBIC     string `json:"receivingBIC"`
		ReceivingAccount string `json:"receivingAccountNumber"`
		ReceivingName    string `json:"receivingName"`
		SenderName       string `json:"senderName"`
		SenderBIC        string `json:"senderBIC"`
		SenderAccount    string `json:"senderAccountNumber"`
		Amount           string `json:"amount"`
		Currency         string `json:"currency"`
		InstructionId    string `json:"instructionId"`
		LocalInstrument  string `json:"localInstrument"`
		// ---
		ReferenceId       string `json:"referenceNumber"`   // 12 Digits
		InstapayReference string `json:"instapayReference"` //
		CoreReference     string `json:"coreReference"`
	}

	JanusTransactCredit struct {
		ReceivingBIC     string `json:"receivingBIC"`
		ReceivingAccount string `json:"receivingAccountNumber"`
		ReceivingName    string `json:"receivingName"`
		SenderName       string `json:"senderName"`
		SenderBIC        string `json:"senderBIC"`
		SenderAccount    string `json:"senderAccountNumber"`
		Amount           string `json:"amount"`
		Currency         string `json:"currency"`
		LocalInstrument  string `json:"localInstrument"`
		// ---
		ReferenceNumber   string `json:"referenceNumber"`   // 12 Digits
		InstapayReference string `json:"instapayReference"` // instruction id
		CoreReference     string `json:"coreReference"`     // galing kay iGate
	}
)
