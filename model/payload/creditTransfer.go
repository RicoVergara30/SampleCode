package payload

// JSON
type (
	FDSAPRequestCreditTransfer struct {
		ReceivingBIC           string  `json:"receivingBIC"`
		ReceivingAccountNumber string  `json:"receivingAccountNumber"`
		ReceivingName          string  `json:"receivingName"`
		SenderBIC              string  `json:"senderBIC"`
		SenderName             string  `json:"senderName"`
		SenderAccountNumber    string  `json:"senderAccountNumber"`
		Amount                 float64 `json:"amount"` // 54
		Currency               string  `json:"currency"`
		LocalInstrument        string  `json:"localInstrument"`
		ReferenceID            string  `json:"referenceId"`
		AppId                  string  `json:"appId,omitempty"`
		Application            string  `json:"application"`
		// IDs
		BusinessMessageId   string `json:"businessMessageId"`
		InstructionId       string `json:"instructionId"`
		TransactionId       string `json:"transactionId"`
		MessageDefenitionId string `json:"messageDefinitionId"`
		// QR Fields
		TypeOfTrans              string `json:"typeOfTrans"`
		PaymentSystemUId         string `json:"paymentSysUId"`                 // 27-00
		AquiredId                string `json:"aquiredId"`                     // 27-01
		PaymentType              string `json:"paymentType"`                   // 27-02
		MechId                   string `json:"merchId"`                       // 27-03
		MerchCreditAccount       string `json:"merchCreditAccount"`            // 27-04
		MobileNumber             string `json:"mobileNumber"`                  // 27-05
		MerchCode                string `json:"merchCode"`                     // 52
		CurrCode                 string `json:"currCode"`                      // 53
		CountryCode              string `json:"countryCode"`                   // 58
		MerchName                string `json:"merchName"`                     // 59
		MerchCity                string `json:"merchCity"`                     // 60
		PostalCode               string `json:"postalCode"`                    // 61
		GlobalUniqueIdentifier   string `json:"globallyUniqueIdentifier"`      // 62-00
		BillNumber               string `json:"billNumber"`                    // 62-01
		AddtlMobileNumber        string `json:"additionalMobileNumber"`        // 64-02
		StoreLabel               string `json:"storeLabel"`                    // 62-03
		LoyaltyNumber            string `json:"loyaltyNumber"`                 // 62-04
		ReferenceLabel           string `json:"referenceLabel"`                // 62-05
		CustomerLabel            string `json:"customerLabel"`                 // 62-06
		TerminalLabel            string `json:"terminalLabel"`                 // 62-07 // Include this if transaction is P2M/P2B
		PurposeTrans             string `json:"purposeTrans"`                  // 62-08
		AddtlConsumerDataRequest string `json:"additionalConsumerDataRequest"` // 62-09
	}

	FDSIRequestCreditTransfer struct {
		SenderAccountNumber    string  `json:"SenderAccountNumber,omitempty"`
		SenderAccountName      string  `json:"SenderAccountName,omitempty"`
		RecipientBankCode      string  `json:"RecipientBankCode,omitempty"`
		RecipientAccountNumber string  `json:"RecipientAccountNumber,omitempty"`
		RecipientAccountName   string  `json:"RecipientAccountName,omitempty"`
		TransactionReference   string  `json:"TransactionReference,omitempty"`
		TransactionAmount      float64 `json:"TransactionAmount,omitempty"`
		TransactionCharge      float64 `json:"TransactionCharge,omitempty"`
		// IDs
		BusinessMessageId int `json:"businessMessageId,omitempty"`
		MessageId         int `json:"messageId,omitempty"`
		InstructionId     int `json:"instructionId,omitempty"`
	}

	ResponseCreditTransfer struct {
		TransactionType     string  `json:"transactionType"`
		Status              string  `json:"status"`
		ReasonCode          string  `json:"reasonCode,omitempty"`
		Description         string  `json:"description,omitempty"`
		LocalInstrument     string  `json:"localInstrument,omitempty"`
		InstructionID       string  `json:"instructionId,omitempty"`
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
		Pacs_008            string  `json:"pacs008,omitempty"`
		Pacs_002            string  `json:"pacs002,omitempty"`
	}

	SystemParameter struct {
		Parameter string `json:"parameter"`
		Value     string `json:"value"`
	}

	ReasonCode struct {
		Code        string `json:"code"`
		Description string `json:"description"`
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
