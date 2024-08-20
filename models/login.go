package models

type (
	LoginPage struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	ResponseLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Status   bool   `jso:"status"`
		Token    string `json:"token"`
	}
	RetCodes struct {
		RetCode string `json:"retcode"`
		Message string `json:"message"`
	}
)
