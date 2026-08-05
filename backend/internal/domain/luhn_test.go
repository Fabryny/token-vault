package domain_test

import (
	"testing"

	"github.com/Fabryny/token-vault/backend/internal/domain"
)

func TestIsValidLuhn(t *testing.T) {
	cases := []struct {
		name string
		pan  string
		want bool
	}{
		{"visa de teste válido", "4111111111111111", true},
		{"mastercard de teste", "5555555555554444", true},
		{"amex de teste (15)", "378282246310005", true},
		{"um dígito trocado", "4111111111111112", false},
		{"vizinhos transpostos", "4111111111111611", false},
		{"string vazia", "", false},
		{"curto demais", "411111111111", false},
		{"longo demais", "41111111111111111111", false},
		{"contém letra", "4111111111111a11", false},
		{"contém sinal", "-411111111111111", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := domain.IsValidLuhn(tc.pan); got != tc.want {
				t.Errorf("IsValidLuhn(%q) = %v, quer %v", tc.pan, got, tc.want)
			}
		})
	}
}
