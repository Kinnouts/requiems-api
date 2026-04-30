package seedutil

import "testing"

func TestToInt16(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		v       int
		field   string
		want    int16
		wantErr bool
	}{
		{name: "zero", v: 0, field: "f", want: 0},
		{name: "positive in range", v: 100, field: "f", want: 100},
		{name: "max int16", v: 32767, field: "f", want: 32767},
		{name: "min int16", v: -32768, field: "f", want: -32768},
		{name: "above max", v: 32768, field: "f", wantErr: true},
		{name: "below min", v: -32769, field: "f", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ToInt16(tt.v, tt.field)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ToInt16(%d) error = %v, wantErr %v", tt.v, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("ToInt16(%d) = %d, want %d", tt.v, got, tt.want)
			}
		})
	}
}
