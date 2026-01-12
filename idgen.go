// Package idgen implements ID generation
//
// Copyright 2026 William J House.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package idgen

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func DocID(prefix string, format string) string {
	return generateDocId(prefix, format)
}

// Helper function to generate document ID for Couchbase
func generateDocId(prefix string, format string) string {
	// Try to get type prefix from the row data or use default

	typePrefix := prefix
	if prefix == "" {
		typePrefix = "doc"
	}

	// Clean the prefix (remove spaces, convert to lowercase)
	typePrefix = strings.ToLower(strings.ReplaceAll(typePrefix, " ", "_"))

	// Check if there's a preference for ID format (could be set in metadata or config)
	// For now, default to the shorter format, but make it configurable
	switch format {
	case "uuid":
		return fmt.Sprintf("%s:%s", typePrefix, uuid.New().String())
	case "timestamp":
		return fmt.Sprintf("%s:%d", typePrefix, time.Now().UnixNano())
	case "short":
		// useShortFormat = true
		s, _ := cryptoRandomString(10)
		return fmt.Sprintf("%s:%s", typePrefix, s)
	case "crypto":
		// Generate a cryptographic random ID
		s, _ := cryptoRandomString(16)
		return fmt.Sprintf("%s:%s", typePrefix, s)
	default:
		// Default to short format
		s, _ := cryptoRandomString(10)
		return fmt.Sprintf("%s:%s", typePrefix, s)
	}
}

// cryptoRandomString(n int) - generates a cryptographically secure random string of length n
func cryptoRandomString(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	out := make([]byte, n)

	// To avoid modulo bias, only accept random bytes < 62*4 = 248.
	// 248 is the largest multiple of 62 less than 256.
	const max = byte(248)

	i := 0
	buf := make([]byte, n*2) // amortize syscalls; ok to overshoot a bit
	for i < n {
		if _, err := rand.Read(buf); err != nil {
			return "", err
		}
		for _, b := range buf {
			if b >= max {
				continue
			}
			out[i] = charset[b%62]
			i++
			if i == n {
				break
			}
		}
	}

	return string(out), nil
}
