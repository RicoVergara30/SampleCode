package webtool

type (
	SettlementFields struct {
		AccountNumber string `json:"accountNumber"`
		Event         string `json:"event"`
		IsEnabled     bool   `json:"isEnabled"`
		Description   string `json:"description"`
	}

	SettlementAccount struct {
		AccountNumber string `json:"accountNumber"`
	}
)

// ENCRYPTION FIELDS
type (
	EncryptionDBUser struct {
		SecretKey string      `json:"secretKey"`
		Data      interface{} `json:""`
		Host      string      `json:"host"`
		Dbname    string      `json:"dbName"`
		Username  string      `json:"username"`
		Password  string      `json:"password"`
	}
)
