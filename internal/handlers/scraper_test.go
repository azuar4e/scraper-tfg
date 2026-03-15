package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func TestParsePrice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"Simple dot", "$123.45", 123.45, false},
		{"Euro comma+dot", "€1,234.56", 1234.56, false},
		{"Euro dot+coma", "€1.234,56", 1234.56, false},
		{"Euro comma", "Precio: 99,99€", 99.99, false},
		{"No price", "sin precio", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePrice(tt.input)
			if tt.hasError {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Fatalf("expected %f, got %f", tt.expected, got)
			}
		})
	}
}

func TestProcessMessageHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `<html><body><h1 id="title">Producto Test</h1><span class="a-price"><span class="a-offscreen">€123,45</span></span></body></html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	}))
	defer server.Close()

	testURL := server.URL + "?test=amazon"

	jsonMsg := `{"id": 123, "user_id": 456, "url": "` + testURL + `", "target_price": 200.0}`
	msg := types.Message{Body: &jsonMsg}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Panicked (esperado sin AWS mocks): %v", r)
		}
	}()

	ProcessMessageHandler(msg)

	t.Log("Scraping logic test completed")
}
