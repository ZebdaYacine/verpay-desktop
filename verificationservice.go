package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
)

// FileVerificationResult is returned after the backend has inspected an S2IR file.
type FileVerificationResult struct {
	FileName    string `json:"fileName"`
	Hash        string `json:"hash"`
	Size        int    `json:"size"`
	Valid       bool   `json:"valid"`
	HashMatches bool   `json:"hashMatches"`
	Message     string `json:"message"`
}

// FileVerificationService receives the selected file from the frontend and
// performs the validation that must not be trusted to the browser alone.
type FileVerificationService struct{}

// VerifyFile validates the file extension, computes its SHA-256 hash and, when
// provided, compares it with the expected hash.
//
// Wails serializes []byte arguments as base64 strings in the generated binding,
// so the frontend can send the selected File contents directly through the RPC.
func (s *FileVerificationService) VerifyFile(
	fileName string,
	contents []byte,
	expectedHash string,
) (FileVerificationResult, error) {

	ext := strings.ToLower(filepath.Ext(fileName))

	if ext != ".s2ir" && ext != ".txt" {
		return FileVerificationResult{},
			fmt.Errorf("invalid file extension: %s", ext)
	}

	hash := computeSHA256(contents)

	if expectedHash != "" &&
		!strings.EqualFold(hash, expectedHash) {
		return FileVerificationResult{}, fmt.Errorf("SHA256 mismatch")
	}

	valid, err := isValidS2IRFile(contents)
	if err != nil {
		return FileVerificationResult{}, err
	}

	return FileVerificationResult{
		FileName: filepath.Base(fileName),
		Hash:     hash,
		Size:     len(contents),
		Valid:    valid,
		Message:  "",
	}, nil
}

func isValidS2IRFile(contents []byte) (bool, error) {

	result, err := ValidateCCPFile(contents)

	if err != nil {
		return false, err
	}

	return result.Valid, nil
}
func (s *FileVerificationService) GenerateOutPutFile(
	fileName string,
	contents []byte,
	expectedHash string,
) (FileVerificationResult, error) {
	return FileVerificationResult{}, nil
}

func computeSHA256(contents []byte) string {
	hash := sha256.Sum256(contents)
	return hex.EncodeToString(hash[:])
}
