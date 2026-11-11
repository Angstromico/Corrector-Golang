package token_test

import (
	"encoding/json"
	"testing"

	"corrector-ortografico/internal/core/token"
)

func TestSpan_Len(t *testing.T) {
	tests := []struct {
		name     string
		span     token.Span
		expected int
	}{
		{"standard span", token.Span{Start: 5, End: 12, Line: 1, Column: 6}, 7},
		{"empty span", token.Span{Start: 3, End: 3, Line: 1, Column: 4}, 0},
		{"inverted span defaults to 0", token.Span{Start: 10, End: 5, Line: 1, Column: 11}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.span.Len(); got != tt.expected {
				t.Errorf("Span.Len() = %d, expected %d", got, tt.expected)
			}
		})
	}
}

func TestSpan_Contains(t *testing.T) {
	span := token.Span{Start: 10, End: 20, Line: 2, Column: 5}

	tests := []struct {
		offset   int
		expected bool
	}{
		{9, false},
		{10, true},
		{15, true},
		{19, true},
		{20, false},
		{25, false},
	}

	for _, tt := range tests {
		if got := span.Contains(tt.offset); got != tt.expected {
			t.Errorf("Span.Contains(%d) = %v, expected %v", tt.offset, got, tt.expected)
		}
	}
}

func TestSpan_Merge(t *testing.T) {
	s1 := token.Span{Start: 10, End: 15, Line: 2, Column: 5}
	s2 := token.Span{Start: 18, End: 25, Line: 2, Column: 13}

	merged1 := s1.Merge(s2)
	if merged1.Start != 10 || merged1.End != 25 || merged1.Line != 2 || merged1.Column != 5 {
		t.Errorf("Merge s1 with s2 = %+v, expected Start 10, End 25, Line 2, Column 5", merged1)
	}

	merged2 := s2.Merge(s1)
	if merged2.Start != 10 || merged2.End != 25 || merged2.Line != 2 || merged2.Column != 5 {
		t.Errorf("Merge s2 with s1 = %+v, expected Start 10, End 25, Line 2, Column 5", merged2)
	}
}

func TestTokenType_String(t *testing.T) {
	tests := []struct {
		tokenType token.TokenType
		expected  string
	}{
		{token.TokenUnknown, "Unknown"},
		{token.TokenWord, "Word"},
		{token.TokenWhitespace, "Whitespace"},
		{token.TokenPunctuation, "Punctuation"},
		{token.TokenInvertedPunctuation, "InvertedPunctuation"},
		{token.TokenNumber, "Number"},
		{token.TokenSymbol, "Symbol"},
		{token.TokenEOF, "EOF"},
		{token.TokenType(99), "TokenType(99)"},
	}

	for _, tt := range tests {
		if got := tt.tokenType.String(); got != tt.expected {
			t.Errorf("TokenType(%d).String() = %q, expected %q", tt.tokenType, got, tt.expected)
		}
	}
}

func TestToken_JSON(t *testing.T) {
	span := token.Span{Start: 0, End: 4, Line: 1, Column: 1}
	tok := token.NewToken(token.TokenWord, "hola", span)
	tok.POSTag = "NOUN"
	tok.Lemma = "hola"
	tok.Feats["gender"] = "fem"

	data, err := json.Marshal(tok)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var unmarshaled token.Token
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if unmarshaled.Value != "hola" || unmarshaled.POSTag != "NOUN" || unmarshaled.Span.Len() != 4 {
		t.Errorf("Unexpected unmarshaled token: %+v", unmarshaled)
	}
}
