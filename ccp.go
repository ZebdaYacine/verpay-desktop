package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const (
	CCPLineLength = 62
)

// ============================================================
// STRUCTURES
// ============================================================


type CCPDuplicateType string

const (
	CCPDuplicateExact CCPDuplicateType = "EXACT"

	CCPDuplicateDifferentName CCPDuplicateType = "DIFFERENT_NAME"
)

type CCPDuplicate struct {
	Type       CCPDuplicateType `json:"type"`
	FirstLine  int              `json:"first_line"`
	SecondLine int              `json:"second_line"`

	CCP        string `json:"ccp"`
	Amount     string `json:"amount"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
}


// ============================================================
// MAIN VALIDATION
// ============================================================

// ValidateCCPFile validates an S2IR file line by line.
func ValidateCCPFile(
	contents []byte,
) (*CCPValidationResult, error) {

	scanner := bufio.NewScanner(bytes.NewReader(contents))

	scanner.Buffer(
		make([]byte, 64*1024),
		1024*1024,
	)

	result := &CCPValidationResult{
		Payments:   make([]CCPPayment, 0),
		Duplicates: make([]CCPDuplicate, 0),
		Valid:      false,
	}

	lineNumber := 0
	totalAmount := new(big.Int)

	// ========================================================
	// READ FILE LINE BY LINE
	// ========================================================

	for scanner.Scan() {

		lineNumber++
		line := scanner.Text()

		// ----------------------------------------------------
		// 1. VALIDATE LINE LENGTH
		// ----------------------------------------------------

		if len(line) != CCPLineLength {
			return nil, fmt.Errorf(
				"line %d: invalid length, expected %d, got %d",
				lineNumber,
				CCPLineLength,
				len(line),
			)
		}

		// ----------------------------------------------------
		// 2. VALIDATE FIRST CHARACTER
		// ----------------------------------------------------

		if line[0] != '*' {
			return nil, fmt.Errorf(
				"line %d: must start with '*'",
				lineNumber,
			)
		}

		// ----------------------------------------------------
		// 3. IDENTIFY RECORD TYPE
		// ----------------------------------------------------

		recordType := line[61]

		switch recordType {

		// ====================================================
		// HEADER - TYPE 0
		// ====================================================

		case '0':

			if lineNumber != 1 {
				return nil, fmt.Errorf(
					"line %d: header must be the first line",
					lineNumber,
				)
			}

			header, err := validateCCPHeader(line)

			if err != nil {
				return nil, fmt.Errorf(
					"line %d: %w",
					lineNumber,
					err,
				)
			}

			result.Header = header

		// ====================================================
		// PAYMENT - TYPE 1
		// ====================================================

		case '1':

			if lineNumber == 1 {
				return nil, fmt.Errorf(
					"line %d: payment found before header",
					lineNumber,
				)
			}

			payment, err := validateCCPPayment(
				line,
				lineNumber,
			)

			if err != nil {
				return nil, fmt.Errorf(
					"line %d: %w",
					lineNumber,
					err,
				)
			}

			// Convert amount to integer.
			amount, ok := new(big.Int).SetString(
				payment.Amount,
				10,
			)

			if !ok {
				return nil, fmt.Errorf(
					"line %d: invalid payment amount",
					lineNumber,
				)
			}

			// Add amount to global total.
			totalAmount.Add(totalAmount, amount)

			// Increase payment count.
			result.PaymentCount++

			// Save payment.
			result.Payments = append(
				result.Payments,
				payment,
			)

		// ====================================================
		// INVALID RECORD TYPE
		// ====================================================

		default:

			return nil, fmt.Errorf(
				"line %d: invalid record type %q",
				lineNumber,
				recordType,
			)
		}
	}

	// ========================================================
	// SCANNER ERRORS
	// ========================================================

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(
			"error reading CCP file: %w",
			err,
		)
	}

	// ========================================================
	// CHECK EMPTY FILE
	// ========================================================

	if lineNumber == 0 {
		return nil, fmt.Errorf("CCP file is empty")
	}

	// ========================================================
	// CHECK PAYMENT COUNT
	// ========================================================

	if result.PaymentCount == 0 {
		return nil, fmt.Errorf(
			"CCP file contains no payments",
		)
	}

	if result.PaymentCount != result.Header.PaymentCount {
		return nil, fmt.Errorf(
			"payment count mismatch: header=%d, body=%d",
			result.Header.PaymentCount,
			result.PaymentCount,
		)
	}

	// ========================================================
	// CHECK HEADER GLOBAL AMOUNT
	// ========================================================

	headerAmount, ok := new(big.Int).SetString(
		result.Header.TotalAmount,
		10,
	)

	if !ok {
		return nil, fmt.Errorf(
			"invalid header global amount",
		)
	}

	// Compare sum of body amounts with header amount.
	if totalAmount.Cmp(headerAmount) != 0 {
		return nil, fmt.Errorf(
			"amount mismatch: header=%s, body=%s",
			headerAmount.String(),
			totalAmount.String(),
		)
	}

	// ========================================================
	// DETECT DUPLICATES
	// ========================================================

	result.Duplicates = DetectCCPDuplicates(
		result.Payments,
	)

	result.DuplicateCount = len(result.Duplicates)

	for _, duplicate := range result.Duplicates {

		switch duplicate.Type {

		case CCPDuplicateExact:
			result.ExactDuplicateCount++

		case CCPDuplicateDifferentName:
			result.DifferentNameDuplicateCount++
		}
	}

	// ========================================================
	// SUCCESS
	// ========================================================

	result.TotalAmount = totalAmount.String()
	result.Valid = true

	if result.DuplicateCount > 0 {

		result.Message = fmt.Sprintf(
			"CCP file valid, but %d duplicate(s) detected",
			result.DuplicateCount,
		)

	} else {

		result.Message = "CCP file verified successfully"
	}

	return result, nil
}

// ============================================================
// HEADER VALIDATION
// ============================================================

func validateCCPHeader(line string) (CCPHeader, error) {

	header := CCPHeader{
		CCP:         strings.TrimSpace(line[1:21]),
		TotalAmount: strings.TrimSpace(line[21:34]),
		Period:      line[41:47],
	}

	paymentCount := strings.TrimSpace(line[34:41])

	// Validate CCP.
	if !isCCPNumeric(header.CCP) {
		return header, fmt.Errorf(
			"invalid header CCP: %q",
			header.CCP,
		)
	}

	// Validate total amount.
	if !isCCPNumeric(header.TotalAmount) {
		return header, fmt.Errorf(
			"invalid header total amount: %q",
			header.TotalAmount,
		)
	}

	// Validate payment count.
	if !isCCPNumeric(paymentCount) {
		return header, fmt.Errorf(
			"invalid header payment count: %q",
			paymentCount,
		)
	}

	count, err := strconv.Atoi(paymentCount)

	if err != nil {
		return header, fmt.Errorf(
			"invalid payment count: %w",
			err,
		)
	}

	header.PaymentCount = count

	// Validate period: MMAAAA.
	if len(header.Period) != 6 ||
		!isCCPNumeric(header.Period) {

		return header, fmt.Errorf(
			"invalid period: %q",
			header.Period,
		)
	}

	month, _ := strconv.Atoi(header.Period[0:2])
	year, _ := strconv.Atoi(header.Period[2:6])

	if month < 1 || month > 12 || year == 0 {
		return header, fmt.Errorf(
			"invalid period: %s",
			header.Period,
		)
	}

	// Validate 14 spaces at positions 48-61.
	if line[47:61] != strings.Repeat(" ", 14) {

		return header, fmt.Errorf(
			"positions 48-61 must contain 14 spaces",
		)
	}

	return header, nil
}

// ============================================================
// PAYMENT VALIDATION
// ============================================================

func validateCCPPayment(
	line string,
	lineNumber int,
) (CCPPayment, error) {

	payment := CCPPayment{
		LineNumber:  lineNumber,
		CCP:         strings.TrimSpace(line[1:21]),
		Amount:      strings.TrimSpace(line[21:34]),
		Beneficiary: strings.TrimSpace(line[34:61]),
	}

	// Validate CCP.
	if !isCCPNumeric(payment.CCP) {
		return payment, fmt.Errorf(
			"invalid beneficiary CCP: %q",
			payment.CCP,
		)
	}

	// Validate amount.
	if !isCCPNumeric(payment.Amount) {
		return payment, fmt.Errorf(
			"invalid payment amount: %q",
			payment.Amount,
		)
	}

	// Validate beneficiary.
	if payment.Beneficiary == "" {
		return payment, fmt.Errorf(
			"beneficiary name cannot be empty",
		)
	}

	return payment, nil
}

// ============================================================
// DUPLICATE DETECTION
// ============================================================

// DetectCCPDuplicates detects:
// 1. Same CCP + same amount + same name.
// 2. Same CCP + same amount + different name.
//
// Each additional occurrence is compared to the first
// occurrence of the same CCP and amount.
func DetectCCPDuplicates(
	payments []CCPPayment,
) []CCPDuplicate {

	duplicates := make([]CCPDuplicate, 0)

	type SeenPayment struct {
		Payment CCPPayment
		Line    int
	}

	// Key: CCP + normalized amount.
	seen := make(map[string]SeenPayment)

	for index, payment := range payments {

		currentLine := payment.LineNumber

		// Fallback for manually constructed payments.
		if currentLine == 0 {
			currentLine = index + 2
		}

		// Normalize amount.
		amount, ok := new(big.Int).SetString(
			strings.TrimSpace(payment.Amount),
			10,
		)

		if !ok {
			continue
		}

		normalizedAmount := amount.String()
		ccp := strings.TrimSpace(payment.CCP)

		key := ccp + "|" + normalizedAmount

		previous, exists := seen[key]

		// First occurrence.
		if !exists {

			seen[key] = SeenPayment{
				Payment: payment,
				Line:    currentLine,
			}

			continue
		}

		// A repeated CCP and amount was found.
		duplicateType := CCPDuplicateExact

		firstName := normalizeCCPName(
			previous.Payment.Beneficiary,
		)

		secondName := normalizeCCPName(
			payment.Beneficiary,
		)

		if firstName != secondName {
			duplicateType = CCPDuplicateDifferentName
		}

		duplicates = append(
			duplicates,
			CCPDuplicate{
				Type:       duplicateType,
				FirstLine:  previous.Line,
				SecondLine: currentLine,
				CCP:        ccp,
				Amount:     normalizedAmount,
				FirstName:  previous.Payment.Beneficiary,
				SecondName: payment.Beneficiary,
			},
		)
	}

	return duplicates
}

// ============================================================
// NAME NORMALIZATION
// ============================================================

// normalizeCCPName ignores:
// - Leading/trailing spaces
// - Repeated spaces
// - Uppercase/lowercase differences
func normalizeCCPName(name string) string {

	return strings.ToUpper(
		strings.Join(
			strings.Fields(name),
			" ",
		),
	)
}

// ============================================================
// NUMERIC VALIDATION
// ============================================================

func isCCPNumeric(value string) bool {

	if value == "" {
		return false
	}

	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// ============================================================
// SERVICE COMPATIBILITY
// ============================================================

