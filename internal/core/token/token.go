package token

// Span defines a 0-indexed character range and 1-indexed line/column position in the text.
type Span struct {
	Start  int `json:"start"`  // Start character offset (0-indexed, inclusive)
	End    int `json:"end"`    // End character offset (0-indexed, exclusive)
	Line   int `json:"line"`   // Line number (1-indexed)
	Column int `json:"column"` // Column number (1-indexed)
}

// Len returns the length of the span in characters.
func (s Span) Len() int {
	if s.End < s.Start {
		return 0
	}
	return s.End - s.Start
}

// Contains reports whether offset is within the span [Start, End).
func (s Span) Contains(offset int) bool {
	return offset >= s.Start && offset < s.End
}

// Merge combines two spans into a single span spanning from the earliest start
// to the latest end. Line and Column are taken from the earlier start span.
func (s Span) Merge(other Span) Span {
	start := s.Start
	line := s.Line
	column := s.Column

	if other.Start < start {
		start = other.Start
		line = other.Line
		column = other.Column
	}

	end := s.End
	if other.End > end {
		end = other.End
	}

	return Span{
		Start:  start,
		End:    end,
		Line:   line,
		Column: column,
	}
}

// Token represents a lexical token and its morphosyntactic metadata.
type Token struct {
	Type   TokenType         `json:"type"`
	Value  string            `json:"value"`
	Span   Span              `json:"span"`
	POSTag string            `json:"pos_tag,omitempty"` // e.g. "NOUN", "DET", "VERB"
	Lemma  string            `json:"lemma,omitempty"`   // e.g. "cantar" for "cantaban"
	Feats  map[string]string `json:"feats,omitempty"`   // e.g. {"gender": "fem", "number": "sing"}
}

// NewToken constructs a new Token with empty Feats map.
func NewToken(tokenType TokenType, value string, span Span) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
		Span:  span,
		Feats: make(map[string]string),
	}
}
