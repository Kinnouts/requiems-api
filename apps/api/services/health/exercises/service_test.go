package exercises

import (
	"testing"
)

func TestExercise_IsData(t *testing.T) {
	var e Exercise
	e.IsData()
}

func TestExerciseList_IsData(t *testing.T) {
	var l ExerciseList
	l.IsData()
}

func TestStringList_IsData(t *testing.T) {
	var s StringList
	s.IsData()
}

func TestExercise_FieldsPresent(t *testing.T) {
	e := Exercise{
		ID:               1,
		Name:             "band shrug",
		BodyParts:        []string{"neck"},
		Equipment:        []string{"band"},
		TargetMuscles:    []string{"traps"},
		SecondaryMuscles: []string{"shoulders"},
		Instructions:     []string{"Stand with feet shoulder-width apart."},
	}

	if e.ID != 1 {
		t.Errorf("expected ID 1, got %d", e.ID)
	}
	if e.Name != "band shrug" {
		t.Errorf("expected name 'band shrug', got %q", e.Name)
	}
	if len(e.BodyParts) != 1 || e.BodyParts[0] != "neck" {
		t.Errorf("unexpected body_parts: %v", e.BodyParts)
	}
	if len(e.Instructions) != 1 {
		t.Errorf("expected 1 instruction, got %d", len(e.Instructions))
	}
}

func TestExerciseList_Pagination(t *testing.T) {
	l := ExerciseList{
		Items:   []Exercise{{ID: 1, Name: "squat"}, {ID: 2, Name: "deadlift"}},
		Total:   50,
		Page:    2,
		PerPage: 20,
	}

	if l.Total != 50 {
		t.Errorf("expected total 50, got %d", l.Total)
	}
	if l.Page != 2 {
		t.Errorf("expected page 2, got %d", l.Page)
	}
	if l.PerPage != 20 {
		t.Errorf("expected per_page 20, got %d", l.PerPage)
	}
	if len(l.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(l.Items))
	}
}

func TestStringList_ItemsAndTotal(t *testing.T) {
	s := StringList{
		Items: []string{"chest", "back", "legs"},
		Total: 3,
	}

	if s.Total != 3 {
		t.Errorf("expected total 3, got %d", s.Total)
	}
	if len(s.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(s.Items))
	}
}

// TestListParams_ZeroValue verifies the zero value of ListParams has empty
// filter fields. The page/per_page defaults (1/20) are applied by the HTTP
// handler before binding, not by the struct itself.
func TestListParams_ZeroValue(t *testing.T) {
	var p ListParams

	if p.BodyPart != "" || p.Equipment != "" || p.Muscle != "" || p.Search != "" {
		t.Error("expected all filter fields to be empty on zero value")
	}
	if p.Page != 0 || p.PerPage != 0 {
		t.Error("expected page and per_page to be zero on zero value")
	}
}
