package main

// CCPHeader represents a header record (type 0).
type CCPHeader struct {
	CCP          string `json:"ccp"`
	TotalAmount  string `json:"total_amount"`
	PaymentCount int    `json:"payment_count"`
	Period       string `json:"period"`
}

// CCPPayment represents a payment record (type 1).
type CCPPayment struct {
	CCP         string `json:"ccp"`
	Amount      string `json:"amount"`
	Beneficiary string `json:"beneficiary"`
}

// CCPValidationResult contains validated file information.
type CCPValidationResult struct {
	Header       CCPHeader    `json:"header"`
	Payments     []CCPPayment `json:"payments"`
	PaymentCount int          `json:"payment_count"`
	TotalAmount  string       `json:"total_amount"`
	Valid        bool         `json:"valid"`
	Message      string       `json:"message"`
}
