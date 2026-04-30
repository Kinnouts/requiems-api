package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// ExerciseRecord holds the normalised data to upsert into the exercises table.
// gifUrl and the raw exerciseId are intentionally omitted from the public API;
// exerciseId is kept only as ExternalID for idempotent upserts.
type ExerciseRecord struct {
	ExternalID       string
	Name             string
	BodyParts        []string
	Equipment        []string
	TargetMuscles    []string
	SecondaryMuscles []string
	Instructions     []string
}

// stepPrefix matches the "Step:N " prefix injected by the source dataset.
var stepPrefix = regexp.MustCompile(`^Step:\d+\s+`)

// loadExercises reads exercises.json from dataDir and returns normalised records.
// gifUrl is never read. External IDs are stored for upsert keying only.
func loadExercises(dataDir string) ([]ExerciseRecord, error) {
	path := filepath.Join(dataDir, "exercises.json")

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	// Decode only the fields we actually need; gifUrl is absent from the struct
	// so the decoder will silently discard it.
	type raw struct {
		ExerciseID       string   `json:"exerciseId"`
		Name             string   `json:"name"`
		BodyParts        []string `json:"bodyParts"`
		Equipments       []string `json:"equipments"`
		TargetMuscles    []string `json:"targetMuscles"`
		SecondaryMuscles []string `json:"secondaryMuscles"`
		Instructions     []string `json:"instructions"`
	}

	var items []raw
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	records := make([]ExerciseRecord, 0, len(items))
	for _, item := range items {
		if item.ExerciseID == "" || item.Name == "" {
			continue
		}

		instructions := make([]string, 0, len(item.Instructions))
		for _, step := range item.Instructions {
			instructions = append(instructions, stepPrefix.ReplaceAllString(step, ""))
		}

		records = append(records, ExerciseRecord{
			ExternalID:       item.ExerciseID,
			Name:             item.Name,
			BodyParts:        normalise(item.BodyParts),
			Equipment:        normalise(item.Equipments),
			TargetMuscles:    normalise(item.TargetMuscles),
			SecondaryMuscles: normalise(item.SecondaryMuscles),
			Instructions:     normalise(instructions),
		})
	}

	return records, nil
}

// normalise returns a non-nil slice, replacing nil with an empty slice.
func normalise(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
