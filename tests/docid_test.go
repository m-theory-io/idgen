package tests

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-theory-io/idgen"
)

const allowedCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestDocIDFormats(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		format     string
		wantPrefix string
		validate   func(*testing.T, string)
	}{
		{
			name:       "ShortFormatSanitizedPrefix",
			prefix:     "Invoice Data",
			format:     "short",
			wantPrefix: "invoice_data",
			validate: func(t *testing.T, suffix string) {
				assertCharsetLen(t, suffix, 10)
			},
		},
		{
			name:       "CryptoFormat",
			prefix:     "user",
			format:     "crypto",
			wantPrefix: "user",
			validate: func(t *testing.T, suffix string) {
				assertCharsetLen(t, suffix, 16)
			},
		},
		{
			name:       "UUIDFormat",
			prefix:     "order",
			format:     "uuid",
			wantPrefix: "order",
			validate: func(t *testing.T, suffix string) {
				if _, err := uuid.Parse(suffix); err != nil {
					t.Fatalf("suffix %q is not a valid UUID: %v", suffix, err)
				}
			},
		},
		{
			name:       "TimestampFormat",
			prefix:     "report",
			format:     "timestamp",
			wantPrefix: "report",
			validate: func(t *testing.T, suffix string) {
				if _, err := strconv.ParseInt(suffix, 10, 64); err != nil {
					t.Fatalf("suffix %q is not a valid timestamp: %v", suffix, err)
				}
			},
		},
		{
			name:       "DefaultFormatFallback",
			prefix:     "device",
			format:     "",
			wantPrefix: "device",
			validate: func(t *testing.T, suffix string) {
				assertCharsetLen(t, suffix, 10)
			},
		},
		{
			name:       "EmptyPrefixDefaultsToDoc",
			prefix:     "",
			format:     "uuid",
			wantPrefix: "doc",
			validate: func(t *testing.T, suffix string) {
				if _, err := uuid.Parse(suffix); err != nil {
					t.Fatalf("suffix %q is not a valid UUID: %v", suffix, err)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			id := idgen.DocID(tc.prefix, tc.format)
			parts := strings.SplitN(id, ":", 2)
			if len(parts) != 2 {
				t.Fatalf("expected id to contain prefix and suffix separated by ':' but got %q", id)
			}
			if parts[0] != tc.wantPrefix {
				t.Fatalf("expected prefix %q, got %q", tc.wantPrefix, parts[0])
			}
			tc.validate(t, parts[1])
		})
	}
}

func TestDocIDTiming(t *testing.T) {
	formats := []struct {
		name   string
		format string
	}{
		{name: "Short", format: "short"},
		{name: "Crypto", format: "crypto"},
		{name: "UUID", format: "uuid"},
		{name: "Timestamp", format: "timestamp"},
		{name: "DefaultFallback", format: ""},
	}

	for _, tc := range formats {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			for i := 0; i < 1000; i++ {
				if id := idgen.DocID("timing", tc.format); id == "" {
					t.Fatal("DocID returned empty string")
				}
			}
			elapsed := time.Since(start)
			t.Logf("generated 1000 %s IDs in %s", tc.name, elapsed)
		})
	}
}

func assertCharsetLen(t *testing.T, suffix string, wantLen int) {
	if len(suffix) != wantLen {
		t.Fatalf("expected suffix length %d, got %d (%q)", wantLen, len(suffix), suffix)
	}
	for _, ch := range suffix {
		if !strings.ContainsRune(allowedCharset, ch) {
			t.Fatalf("suffix %q contains invalid character %q", suffix, ch)
		}
	}
}
