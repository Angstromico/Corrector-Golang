package diagnostic_test

import (
	"encoding/json"
	"strings"
	"testing"

	"corrector-ortografico/internal/core/diagnostic"
	"corrector-ortografico/internal/core/token"
)

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		severity diagnostic.Severity
		expected string
	}{
		{diagnostic.SeverityError, "Error"},
		{diagnostic.SeverityWarning, "Warning"},
		{diagnostic.SeverityInfo, "Info"},
		{diagnostic.SeverityHint, "Hint"},
		{diagnostic.Severity(99), "Severity(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.severity.String(); got != tt.expected {
				t.Errorf("Severity(%d).String() = %q, expected %q", tt.severity, got, tt.expected)
			}
		})
	}
}

func TestNewDiagnostic(t *testing.T) {
	span := token.Span{Start: 5, End: 10, Line: 1, Column: 6}
	diag := diagnostic.NewDiagnostic("diag-1", "SPELL001", diagnostic.SeverityError, span, "Palabra desconocida")

	if diag.ID != "diag-1" {
		t.Errorf("ID = %q, expected %q", diag.ID, "diag-1")
	}
	if diag.RuleID != "SPELL001" {
		t.Errorf("RuleID = %q, expected %q", diag.RuleID, "SPELL001")
	}
	if diag.Severity != diagnostic.SeverityError {
		t.Errorf("Severity = %v, expected %v", diag.Severity, diagnostic.SeverityError)
	}
	if diag.Span != span {
		t.Errorf("Span = %+v, expected %+v", diag.Span, span)
	}
	if diag.Message != "Palabra desconocida" {
		t.Errorf("Message = %q, expected %q", diag.Message, "Palabra desconocida")
	}
	if diag.Explanation != "" {
		t.Errorf("Explanation = %q, expected empty string", diag.Explanation)
	}
	if diag.Suggestions == nil || len(diag.Suggestions) != 0 {
		t.Errorf("Suggestions should be initialized and empty, got %+v", diag.Suggestions)
	}
}

func TestDiagnostic_FluentBuilders(t *testing.T) {
	span := token.Span{Start: 0, End: 4, Line: 1, Column: 1}
	diag := diagnostic.NewDiagnostic("diag-2", "AGREE001", diagnostic.SeverityWarning, span, "Discordancia de género").
		WithExplanation("El artículo 'la' no concuerda en género con 'problema'.").
		AddSuggestion(diagnostic.Suggestion{
			Replacement: "el problema",
			Description: "Cambiar artículo a masculino",
			Confidence:  0.95,
		})

	if diag.Explanation != "El artículo 'la' no concuerda en género con 'problema'." {
		t.Errorf("Unexpected Explanation: %q", diag.Explanation)
	}
	if len(diag.Suggestions) != 1 {
		t.Fatalf("Expected 1 suggestion, got %d", len(diag.Suggestions))
	}
	s := diag.Suggestions[0]
	if s.Replacement != "el problema" || s.Description != "Cambiar artículo a masculino" || s.Confidence != 0.95 {
		t.Errorf("Unexpected Suggestion: %+v", s)
	}
}

func TestDiagnostic_JSONSerialization(t *testing.T) {
	span := token.Span{Start: 10, End: 15, Line: 2, Column: 3}
	diag := diagnostic.NewDiagnostic("d-1", "SPELL_01", diagnostic.SeverityError, span, "Error ortográfico").
		AddSuggestion(diagnostic.Suggestion{
			Replacement: "hablar",
			Confidence:  0.99,
		})

	data, err := json.Marshal(diag)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)

	// Verify empty explanation is omitted
	if strings.Contains(jsonStr, `"explanation"`) {
		t.Errorf("Expected 'explanation' to be omitted when empty, got: %s", jsonStr)
	}

	// Verify suggestion with empty description omits description
	if strings.Contains(jsonStr, `"description"`) {
		t.Errorf("Expected 'description' to be omitted when empty in suggestion, got: %s", jsonStr)
	}

	// Verify suggestions are marshaled as array
	if !strings.Contains(jsonStr, `"suggestions":[{"replacement":"hablar","confidence":0.99}]`) {
		t.Errorf("Expected serialized suggestions array, got: %s", jsonStr)
	}

	// Verify unmarshaling round-trip
	var unmarshaled diagnostic.Diagnostic
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if unmarshaled.ID != diag.ID ||
		unmarshaled.RuleID != diag.RuleID ||
		unmarshaled.Severity != diag.Severity ||
		unmarshaled.Span != diag.Span ||
		unmarshaled.Message != diag.Message ||
		len(unmarshaled.Suggestions) != 1 ||
		unmarshaled.Suggestions[0].Replacement != "hablar" {
		t.Errorf("Unmarshaled diagnostic does not match original: %+v", unmarshaled)
	}
}

func TestDiagnostic_JSONEmptySuggestions(t *testing.T) {
	span := token.Span{Start: 0, End: 5, Line: 1, Column: 1}
	diag := diagnostic.NewDiagnostic("d-empty", "INFO_01", diagnostic.SeverityInfo, span, "Información")

	data, err := json.Marshal(diag)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	// Ensure suggestions serializes as [] rather than null for frontend consumption
	if !strings.Contains(string(data), `"suggestions":[]`) {
		t.Errorf("Expected '\"suggestions\":[]' in JSON output, got: %s", string(data))
	}
}

func TestBySpan_Sort(t *testing.T) {
	d1 := *diagnostic.NewDiagnostic("d1", "R1", diagnostic.SeverityError, token.Span{Start: 20, End: 25}, "Third")
	d2 := *diagnostic.NewDiagnostic("d2", "R2", diagnostic.SeverityWarning, token.Span{Start: 5, End: 15}, "Second")
	d3 := *diagnostic.NewDiagnostic("d3", "R3", diagnostic.SeverityInfo, token.Span{Start: 5, End: 10}, "First")
	d4 := *diagnostic.NewDiagnostic("d4", "R4", diagnostic.SeverityError, token.Span{Start: 5, End: 10}, "Zero (higher severity)")

	diags := []diagnostic.Diagnostic{d1, d2, d3, d4}
	diagnostic.Sort(diags)

	// Expected order:
	// 1. Start 5, End 10, Severity 1 (d4)
	// 2. Start 5, End 10, Severity 3 (d3)
	// 3. Start 5, End 15, Severity 2 (d2)
	// 4. Start 20, End 25, Severity 1 (d1)
	expectedIDs := []string{"d4", "d3", "d2", "d1"}
	for i, expectedID := range expectedIDs {
		if diags[i].ID != expectedID {
			t.Errorf("Index %d: expected ID %q, got %q", i, expectedID, diags[i].ID)
		}
	}
}

func TestBySpanPtr_Sort(t *testing.T) {
	d1 := diagnostic.NewDiagnostic("d1", "R1", diagnostic.SeverityError, token.Span{Start: 30, End: 35}, "C")
	d2 := diagnostic.NewDiagnostic("d2", "R2", diagnostic.SeverityWarning, token.Span{Start: 10, End: 20}, "B")
	d3 := diagnostic.NewDiagnostic("d3", "R3", diagnostic.SeverityInfo, token.Span{Start: 10, End: 15}, "A")

	diags := []*diagnostic.Diagnostic{d1, d2, d3}
	diagnostic.SortPtr(diags)

	expectedIDs := []string{"d3", "d2", "d1"}
	for i, expectedID := range expectedIDs {
		if diags[i].ID != expectedID {
			t.Errorf("Index %d: expected ID %q, got %q", i, expectedID, diags[i].ID)
		}
	}
}
