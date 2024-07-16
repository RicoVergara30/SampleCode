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

		TransactionType string `json:"transactiontype"`
		ReasonCode      string `json:"reasoncode"`
		Status          bool   `json:"status"`
		Description     string `json:"description"`
	}
)
