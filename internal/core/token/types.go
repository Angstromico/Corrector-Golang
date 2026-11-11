package token

import "fmt"

// TokenType represents the classification of a lexical token.
type TokenType int

const (
	TokenUnknown TokenType = iota
	TokenWord
	TokenWhitespace
	TokenPunctuation
	TokenInvertedPunctuation // ¿ ¡
	TokenNumber
	TokenSymbol
	TokenEOF
)

// String returns a human-readable representation of TokenType.
func (t TokenType) String() string {
	switch t {
	case TokenUnknown:
		return "Unknown"
	case TokenWord:
		return "Word"
	case TokenWhitespace:
		return "Whitespace"
	case TokenPunctuation:
		return "Punctuation"
	case TokenInvertedPunctuation:
		return "InvertedPunctuation"
	case TokenNumber:
		return "Number"
	case TokenSymbol:
		return "Symbol"
	case TokenEOF:
		return "EOF"
	default:
		return fmt.Sprintf("TokenType(%d)", t)
	}
}
