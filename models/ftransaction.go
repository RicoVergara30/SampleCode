package models

const (
	RBIBankIdentifierCode = "CAMZPHM2XXX"
)

type (
	Ftransaction struct {
		SenderAccountNumber    string  `json:"SenderAccountNumber,omitempty"`
		SenderAccountName      string  `json:"SenderAccountName,omitempty"`
		RecipientBankCode      string  `json:"RecipientBankCode,omitempty"`
		RecipientAccountNumber string  `json:"RecipientAccountNumber,omitempty"`
		RecipientAccountName   string  `json:"RecipientAccountName,omitempty"`
		TransactionReference   string  `json:"TransactionReference,omitempty"`
		TransactionAmount      float64 `json:"TransactionAmount,omitempty"`
		TransactionCharge      float64 `json:"TransactionCharge,omitempty"`

		//Other response
		TransactionType string `json:"transactiontype"`
		InstructionId   string `json:"instructionId"`
		ReasonCode      string `json:"reasoncode"`
		Status          bool   `json:"status"`
		Description     string `json:"description"`
		Application     string `json:"application"`
	}

	ResponseFtransaction struct {
		TransactionType     string  `json:"transactionType"`      //
		Status              string  `json:"status"`               //
		ReasonCode          string  `json:"reasonCode,omitempty"` //
		Description         string  `json:"description,omitempty"`
		LocalInstrument     string  `json:"localInstrument,omitempty"` //
		InstructionID       string  `json:"instructionId,omitempty"`   //
		TransactionID       string  `json:"transactionId,omitempty"`
		ReferenceID         string  `json:"referenceId,omitempty"`
		SenderBIC           string  `json:"senderBIC,omitempty"`
		SenderName          string  `json:"senderName,omitempty"`
		SenderAccount       string  `json:"senderAccount,omitempty"`
		AmountCurrency      string  `json:"amountCurrency,omitempty"`
		SenderAmount        float64 `json:"senderAmount,omitempty"`
		ReceivingBIC        string  `json:"receivingBIC,omitempty"`
		ReceivingName       string  `json:"receivingName,omitempty"`
		ReceivingAccount    string  `json:"receivingAccount,omitempty"`
		TransactionDateTime string  `json:"transactionDateTime,omitempty"`
	}

	Response struct {
		Device      string      `json:"device"`
		RetCode     string      `json:"retCode"`
		Description string      `json:"description"`
		Response    IPSResponse `json:"response"`
		Error       interface{} `json:"errorData,omitempty"`
	}

	IPSResponse struct {
		InstructionID    string `json:"instructionId"`
		TransactionType  string `json:"transactionType,omitempty"`
		Status           string `json:"status,omitempty"`
		ReasonCode       string `json:"reasonCode,omitempty"`
		Description      string `json:"description,omitempty"`
		ReferenceId      string `json:"referenceId,omitempty"`
		SenderBIC        string `json:"senderBIC,omitempty"`
		SenderName       string `json:"senderName,omitempty"`
		SenderAccount    string `json:"senderAccount,omitempty"`
		ReceivingBIC     string `json:"receivingBIC,omitempty"`
		ReceivingName    string `json:"receivingName,omitempty"`
		ReceivingAccount string `json:"receivingAccount,omitempty"`
	}
	ResponseError struct {
		Code        string `json:"reasonCode"`
		Description string `json:"description"`
		Service     string `json:"service"`
	}
)
