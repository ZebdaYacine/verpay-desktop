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
// MAIN VALIDATION
// ============================================================

// ValidateCCPFile validates an entire S2IR file line by line.
func ValidateCCPFile(contents []byte) (*CCPValidationResult, error) {

	scanner := bufio.NewScanner(bytes.NewReader(contents))

	scanner.Buffer(
		make([]byte, 64*1024),
		1024*1024,
	)

	result := &CCPValidationResult{
		Payments: make([]CCPPayment, 0),
		Valid:    false,
	}

	lineNumber := 0

	// Use big.Int to avoid financial precision errors.
	totalAmount := new(big.Int)

	// ========================================================
	// READ FILE LINE BY LINE
	// ========================================================

	for scanner.Scan() {

		lineNumber++

		line := scanner.Text()

		// ----------------------------------------------------
		// 1. CHECK LINE LENGTH
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
		// 2. CHECK LINE START
		// ----------------------------------------------------

		if line[0] != '*' {
			return nil, fmt.Errorf(
				"line %d: line must start with '*'",
				lineNumber,
			)
		}

		// ----------------------------------------------------
		// 3. IDENTIFY RECORD TYPE
		// ----------------------------------------------------

		recordType := line[61]

		switch recordType {

		// ====================================================
		// HEADER TYPE 0
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
		// BODY TYPE 1
		// ====================================================

		case '1':

			if lineNumber == 1 {
				return nil, fmt.Errorf(
					"line %d: payment found before header",
					lineNumber,
				)
			}

			payment, err := validateCCPPayment(line)

			if err != nil {
				return nil, fmt.Errorf(
					"line %d: %w",
					lineNumber,
					err,
				)
			}

			// Convert payment amount to big.Int
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

			// Add payment amount to total
			totalAmount.Add(totalAmount, amount)

			// Increment payment count
			result.PaymentCount++

			// Save payment
			result.Payments = append(
				result.Payments,
				payment,
			)

		// ====================================================
		// UNKNOWN RECORD TYPE
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
	// CHECK SCANNER ERROR
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
			"CCP file does not contain any payments",
		)
	}

	if result.PaymentCount != result.Header.PaymentCount {
		return nil, fmt.Errorf(
			"CCP payment count mismatch: header=%d, body=%d",
			result.Header.PaymentCount,
			result.PaymentCount,
		)
	}

	// ========================================================
	// CHECK GLOBAL AMOUNT
	// ========================================================

	headerAmount, ok := new(big.Int).SetString(
		result.Header.TotalAmount,
		10,
	)

	if !ok {
		return nil, fmt.Errorf(
			"invalid global amount in header",
		)
	}

	// Header total must equal sum of all body amounts.
	if totalAmount.Cmp(headerAmount) != 0 {
		return nil, fmt.Errorf(
			"CCP amount mismatch: header=%s, body=%s",
			headerAmount.String(),
			totalAmount.String(),
		)
	}

	// ========================================================
	// SUCCESS
	// ========================================================

	result.TotalAmount = totalAmount.String()
	result.Valid = true
	result.Message = "CCP file verified successfully"

	return result, nil
}

// ============================================================
// HEADER VALIDATION
// ============================================================

// validateCCPHeader validates record type 0.
func validateCCPHeader(line string) (CCPHeader, error) {

	header := CCPHeader{
		CCP:         strings.TrimSpace(line[1:21]),
		TotalAmount: strings.TrimSpace(line[21:34]),
		Period:      line[41:47],
	}

	paymentCount := strings.TrimSpace(line[34:41])

	// Validate CCP
	if !isCCPNumeric(header.CCP) {
		return header, fmt.Errorf(
			"invalid header CCP: %q",
			header.CCP,
		)
	}

	// Validate total amount
	if !isCCPNumeric(header.TotalAmount) {
		return header, fmt.Errorf(
			"invalid header total amount: %q",
			header.TotalAmount,
		)
	}

	// Validate payment count
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

	// Validate period MMAAAA
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

	// Validate spaces positions 48-61
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

// validateCCPPayment validates record type 1.
func validateCCPPayment(line string) (CCPPayment, error) {

	payment := CCPPayment{
		CCP:         strings.TrimSpace(line[1:21]),
		Amount:      strings.TrimSpace(line[21:34]),
		Beneficiary: strings.TrimSpace(line[34:61]),
	}

	// Validate beneficiary CCP
	if !isCCPNumeric(payment.CCP) {
		return payment, fmt.Errorf(
			"invalid beneficiary CCP: %q",
			payment.CCP,
		)
	}

	// Validate payment amount
	if !isCCPNumeric(payment.Amount) {
		return payment, fmt.Errorf(
			"invalid payment amount: %q",
			payment.Amount,
		)
	}

	// Validate beneficiary name
	if payment.Beneficiary == "" {
		return payment, fmt.Errorf(
			"beneficiary name cannot be empty",
		)
	}

	return payment, nil
}

// ============================================================
// NUMERIC VALIDATION
// ============================================================

// isCCPNumeric checks whether a string contains only digits.
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


