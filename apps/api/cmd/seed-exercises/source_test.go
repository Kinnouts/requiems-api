package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadExercises_HappyPath(t *testing.T) {
	t.Parallel()

	exercises := []map[string]interface{}{
		{
			"exerciseId":       "ex001",
			"name":             "Push-up",
			"bodyParts":        []string{"chest"},
			"equipments":       []string{"body weight"},
			"targetMuscles":    []string{"pectorals"},
			"secondaryMuscles": []string{"triceps"},
			"instructions":     []string{"Step:1 Get into position.", "Step:2 Lower your body."},
			"gifUrl":           "https://example.com/push-up.gif", // must be silently discarded
		},
	}

	dir := writeExercisesJSON(t, exercises)
	records, err := loadExercises(dir)
	if err != nil {
		t.Fatalf("loadExercises: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	r := records[0]
	if r.ExternalID != "ex001" {
		t.Errorf("ExternalID = %q, want %q", r.ExternalID, "ex001")
	}
	if r.Name != "Push-up" {
		t.Errorf("Name = %q, want %q", r.Name, "Push-up")
	}
	// Step prefix should be stripped from instructions.
	if len(r.Instructions) != 2 || r.Instructions[0] != "Get into position." {
		t.Errorf("Instructions[0] = %q, want %q", r.Instructions[0], "Get into position.")
	}
}

func TestLoadExercises_SkipsMissingIDOrName(t *testing.T) {
	t.Parallel()

	exercises := []map[string]interface{}{
		{"exerciseId": "", "name": "No ID Exercise"},               // missing ID
		{"exerciseId": "ex002", "name": ""},                        // missing name
		{"exerciseId": "ex003", "name": "Valid", "bodyParts": nil}, // valid — nil fields normalised
	}

	dir := writeExercisesJSON(t, exercises)
	records, err := loadExercises(dir)
	if err != nil {
		t.Fatalf("loadExercises: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 valid record, got %d", len(records))
	}
	if records[0].ExternalID != "ex003" {
		t.Errorf("ExternalID = %q, want ex003", records[0].ExternalID)
	}
}

func TestLoadExercises_NilSlicesNormalised(t *testing.T) {
	t.Parallel()

	exercises := []map[string]interface{}{
		{"exerciseId": "ex004", "name": "Squat"}, // omit all slice fields
	}

	dir := writeExercisesJSON(t, exercises)
	records, err := loadExercises(dir)
	if err != nil {
		t.Fatalf("loadExercises: %v", err)
	}
	r := records[0]

	// All slice fields should be non-nil empty slices, not nil.
	if r.BodyParts == nil {
		t.Error("BodyParts should not be nil")
	}
	if r.Equipment == nil {
		t.Error("Equipment should not be nil")
	}
	if r.TargetMuscles == nil {
		t.Error("TargetMuscles should not be nil")
	}
	if r.SecondaryMuscles == nil {
		t.Error("SecondaryMuscles should not be nil")
	}
}

func TestLoadExercises_MissingFile(t *testing.T) {
	t.Parallel()

	_, err := loadExercises(t.TempDir())
	if err == nil {
		t.Fatal("expected error for missing exercises.json, got nil")
	}
}

func TestLoadExercises_InvalidJSON(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "exercises.json"), []byte("not-json"), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := loadExercises(dir)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

// writeExercisesJSON marshals exercises to exercises.json in a temp dir and
// returns the dir path.
func writeExercisesJSON(t *testing.T, exercises interface{}) string {
	t.Helper()

	dir := t.TempDir()
	data, err := json.Marshal(exercises)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if err = os.WriteFile(filepath.Join(dir, "exercises.json"), data, 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return dir
}
