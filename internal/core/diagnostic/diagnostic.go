package diagnostic

import (
	"fmt"
	"sort"

	"corrector-ortografico/internal/core/token"
)

// Severity indicates the seriousness and visual treatment of a diagnostic.
type Severity int

const (
	// SeverityError indicates a hard error, e.g., spelling mistake (red wavy underline).
	SeverityError Severity = 1

	// SeverityWarning indicates a potential issue, e.g., grammatical or agreement discordance (yellow/amber wavy underline).
	SeverityWarning Severity = 2

	// SeverityInfo indicates stylistic or typographical guidance (blue wavy underline).
	SeverityInfo Severity = 3

	// SeverityHint represents a minor suggestion or tip.
	SeverityHint Severity = 4
)

// String returns a human-readable representation of Severity.
func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "Error"
	case SeverityWarning:
		return "Warning"
	case SeverityInfo:
		return "Info"
	case SeverityHint:
		return "Hint"
	default:
		return fmt.Sprintf("Severity(%d)", s)
	}
}

// Suggestion represents an actionable fix for a diagnostic.
type Suggestion struct {
	Replacement string  `json:"replacement"`
	Description string  `json:"description,omitempty"`
	Confidence  float64 `json:"confidence"` // 0.0 to 1.0
}

// Diagnostic represents an issue identified in the text along with context and suggestions.
type Diagnostic struct {
	ID          string       `json:"id"`
	RuleID      string       `json:"rule_id"`
	Severity    Severity     `json:"severity"`
	Span        token.Span   `json:"span"`
	Message     string       `json:"message"`
	Explanation string       `json:"explanation,omitempty"`
	Suggestions []Suggestion `json:"suggestions"`
}

// NewDiagnostic creates a new Diagnostic instance initialized with an empty suggestions slice.
func NewDiagnostic(id, ruleID string, severity Severity, span token.Span, message string) *Diagnostic {
	return &Diagnostic{
		ID:          id,
		RuleID:      ruleID,
		Severity:    severity,
		Span:        span,
		Message:     message,
		Suggestions: make([]Suggestion, 0),
	}
}

// WithExplanation sets the detailed explanation for the diagnostic and returns the modified instance.
func (d *Diagnostic) WithExplanation(explanation string) *Diagnostic {
	d.Explanation = explanation
	return d
}

// AddSuggestion appends an actionable suggestion to the diagnostic and returns the modified instance.
func (d *Diagnostic) AddSuggestion(suggestion Suggestion) *Diagnostic {
	d.Suggestions = append(d.Suggestions, suggestion)
	return d
}

// BySpan implements sort.Interface for []Diagnostic based on Start position, End position, and Severity.
type BySpan []Diagnostic

func (a BySpan) Len() int      { return len(a) }
func (a BySpan) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySpan) Less(i, j int) bool {
	if a[i].Span.Start != a[j].Span.Start {
		return a[i].Span.Start < a[j].Span.Start
	}
	if a[i].Span.End != a[j].Span.End {
		return a[i].Span.End < a[j].Span.End
	}
	return a[i].Severity < a[j].Severity
}

// BySpanPtr implements sort.Interface for []*Diagnostic based on Start position, End position, and Severity.
type BySpanPtr []*Diagnostic

func (a BySpanPtr) Len() int      { return len(a) }
func (a BySpanPtr) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a BySpanPtr) Less(i, j int) bool {
	if a[i].Span.Start != a[j].Span.Start {
		return a[i].Span.Start < a[j].Span.Start
	}
	if a[i].Span.End != a[j].Span.End {
		return a[i].Span.End < a[j].Span.End
	}
	return a[i].Severity < a[j].Severity
}

// Sort sorts a slice of Diagnostic values by appearance in text.
func Sort(diagnostics []Diagnostic) {
	sort.Sort(BySpan(diagnostics))
}

// SortPtr sorts a slice of Diagnostic pointers by appearance in text.
func SortPtr(diagnostics []*Diagnostic) {
	sort.Sort(BySpanPtr(diagnostics))
}
