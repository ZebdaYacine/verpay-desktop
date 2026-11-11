package main

type CCPHeader struct {
	CCP          string `json:"ccp"`
	TotalAmount  string `json:"total_amount"`
	PaymentCount int    `json:"payment_count"`
	Period       string `json:"period"`
}

type CCPPayment struct {
	LineNumber  int    `json:"line_number"`
	CCP         string `json:"ccp"`
	Amount      string `json:"amount"`
	Beneficiary string `json:"beneficiary"`
}

type CCPValidationResult struct {
	Header CCPHeader `json:"header"`

	Payments []CCPPayment `json:"payments"`

	PaymentCount int    `json:"payment_count"`
	TotalAmount  string `json:"total_amount"`

	Duplicates     []CCPDuplicate `json:"duplicates"`
	DuplicateCount int            `json:"duplicate_count"`

	ExactDuplicateCount         int `json:"exact_duplicate_count"`
	DifferentNameDuplicateCount int `json:"different_name_duplicate_count"`

	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}


