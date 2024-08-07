package models

type (
	Registration struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	ResponseRes struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)
