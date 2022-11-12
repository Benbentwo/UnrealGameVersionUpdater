package utils

import (
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColorNameValues(t *testing.T) {
	colorMap = map[string]color.Attribute{
		// formatting
		"bold":   color.Bold,
		"faint":  color.Faint,
		"italic": color.Italic,
	}

	colors := ColorNameValues()
	assert.Equal(t, colors, []string{"bold", "faint", "italic"})

}

func TestGetColor(t *testing.T) {
	colorMap = map[string]color.Attribute{
		// formatting
		"bold":   color.Bold,
		"faint":  color.Faint,
		"italic": color.Italic,
	}

	type args struct {
		optionName string
		colorNames []string
	}
	tests := []struct {
		name       string
		args       args
		want       *color.Color
		wantErr    bool
		wantErrStr string
	}{
		{"No Error", args{"bold", []string{"bold"}}, color.New(color.Bold), false, ""},
		{"Error", args{"abc", []string{"abc"}}, nil, true, "Invalid Option abc, abc, [bold faint italic]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetColor(tt.args.optionName, tt.args.colorNames)
			if tt.wantErr {
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
