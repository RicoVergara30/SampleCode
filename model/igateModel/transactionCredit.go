package igatemodel

type (
	RequestTransferCredit struct {
		TransactionType  string  `json:"transactionType"`
		InstructionID    string  `json:"instructionID,omitempty"`
		ReferenceNumber  string  `json:"referenceNumber,omitempty"`
		CreditAccount    string  `json:"creditAccount,omitempty"`
		DebitAccount     string  `json:"debitAccount,omitempty"`
		TransactionFee   float64 `json:"transactionFee,omitempty"`
		SourceBranchCode string  `json:"sourceBranchCode,omitempty"`
		Amount           float64 `json:"amount,omitempty"`
		AdminFee         float64 `json:"adminFee,omitempty"`
		Description      string  `json:"description,omitempty"`

		// for receiving
		ReceivingBIC           string `json:"receivingBIC,omitempty"`
		ReceivingAccountNumber string `json:"receivingAccountNumber,omitempty"`
		ReceivingName          string `json:"receivingName,omitempty"`
		SenderName             string `json:"senderName,omitempty"`
		Currency               string `json:"currency,omitempty"`
		SenderBIC              string `json:"senderBIC,omitempty"`
		SenderAccountNumber    string `json:"senderAccountNumber,omitempty"`
		InstapayReference      string `json:"instapayReference,omitempty"`
		CoreReference          string `json:"coreReference,omitempty"`
		LocalInstrument        string `json:"localInstrument,omitempty"`
	}

	TransferCredit struct {
		ReferenceNumber string  `json:"referenceNumber,omitempty"`
		CreditAccount   string  `json:"creditAccount,omitempty"`
		DebitAccount    string  `json:"debitAccount,omitempty"`
		TransactionFee  float64 `json:"transactionFee,omitempty"`
		Amount          float64 `json:"amount,omitempty"`
		AdminFee        float64 `json:"adminFee,omitempty"`
		Description     string  `json:"description,omitempty"`
	}

	IBFTransferCredit struct {
		ReferenceNumber     string  `json:"referenceNumber"`
		CreditAccount       string  `json:"creditAccount"`
		SendingBankCode     string  `json:"senderBankCode"`
		SendingBankName     string  `json:"senderBankName"`
		SenderAccountNumber string  `json:"senderAccountNumber"`
		IBFTReference       string  `json:"ibftReference"`
		Amount              float64 `json:"amount"`
		AdminFee            float64 `json:"adminFee,omitempty"`
		Description         string  `json:"description,omitempty"`
	}

	// TransferCreditResponse struct {
	// 	ResponseCode           string  `json:"responseCode,omitempty"`
	// 	Description            string  `json:"description,omitempty"`
	// 	CreditAccount          string  `json:"creditAccount,omitempty"`
	// 	DebitAccount           string  `json:"debitAccount,omitempty"`
	// 	CustomerName           string  `json:"customerName,omitempty"`
	// 	AccountName            string  `json:"accountName,omitempty"`
	// 	ReferenceNumber        string  `json:"referenceNumber,omitempty"`
	// 	Amount                 float64 `json:"amount,omitempty"`
	// 	AdminFee               string  `json:"adminFee,omitempty"`
	// 	Reff                   string  `json:"reff,omitempty"`
	// 	CoreReference          string  `json:"coreReference,omitempty"`
	// 	SourceBranchCode       string  `json:"sourceBranchCode,omitempty"`
	// 	DestinationBranchCode  string  `json:"destinationBranchCode,omitempty"`
	// 	SourceProductCode      string  `json:"sourceProductCode,omitempty"`
	// 	DestinationProductCode string  `json:"destinationProductCode,omitempty"`
	// 	DebitCurrency          string  `json:"debitCurrency,omitempty"`
	// 	CreditCurrency         string  `json:"creditCurrency,omitempty"`
	// 	AvalableBalance        string  `json:"availableBalance,omitempty"`
	// 	ArNumber               string  `json:"arNumber,omitempty"`
	// }

	TransferCreditResponse struct {
		ResponseCode       string  `json:"responseCode,omitempty"`
		ReferenceNumber    string  `json:"referenceNumber,omitempty"`
		CreditAccount      string  `json:"creditAccount,omitempty"`
		CreditAccountName  string  `json:"creditAccountName,omitempty"`
		DebitAccount       string  `json:"debitAccount,omitempty"`
		Amount             float64 `json:"amount,omitempty"`
		AdminFee           float64 `json:"adminFee,omitempty"`
		TransReference     string  `json:"transReference,omitempty"`
		CoreReference      string  `json:"coreReference,omitempty"`
		CreditBranchCode   string  `json:"creditBranchCode,omitempty"`
		CreditAvailBalance float64 `json:"creditAvailBalance,omitempty"`
		CreditNarrative    string  `json:"creditNarrative,omitempty"`
		ArNumber           string  `json:"arNumber,omitempty"`
	}
)
