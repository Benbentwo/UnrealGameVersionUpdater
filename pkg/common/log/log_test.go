package log

import (
	"bytes"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestCaptureOutput(t *testing.T) {
	type args struct {
		f func()
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{func() { Logger().Print("") }}, "INFO: \n"},
		{"nonempty", args{func() { Logger().Print("abc") }}, "INFO: abc\n"},
		{"err", args{func() { Logger().Error("err") }}, "ERROR: err\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CaptureOutput(tt.args.f); got != tt.want {
				t.Errorf("CaptureOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLevels(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{"Test", []string{"panic", "fatal", "error", "warning", "info", "debug", "trace"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLevels(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLevels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoggerOK(t *testing.T) {
	labelsPath = t.TempDir()
	got := Logger()
	assert.NotNil(t, got)
}

func TestLoggerFail(t *testing.T) {
	initLogger = func() error {
		return errors.New("mock logger error")
	}
	t.Cleanup(func() {
		initLogger = initializeLogger
	})
	logs := CaptureOutput(func() {
		Logger()
		//assert.Nil(t, got)
	})
	assert.Equal(t, "WARNING: error initializing logrus mock logger error\n", logs)

}

func TestNewVgsTextFormat(t *testing.T) {
	tests := []struct {
		name string
		want *VgsTextFormat
	}{
		{"The Test", &VgsTextFormat{
			ShowInfoLevel:   false,
			ShowTimestamp:   false,
			TimestampFormat: "2006-01-02 15:04:05",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVgsTextFormat(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVgsTextFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetLevel(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"debug", args{"debug"}, false},
		{"Debug", args{"Debug"}, false},
		{"Error", args{"error"}, false},
		{"burrito", args{"burrito"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetLevel(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("SetLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetOutput(t *testing.T) {
	tests := []struct {
		name    string
		wantOut string
	}{
		{"Test", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			SetOutput(out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("SetOutput() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestVgsTextFormat_Format(t *testing.T) {
	var dateFormatString = "2006-01-02"
	currTime := time.Now().Format(dateFormatString)
	type fields struct {
		ShowInfoLevel   bool
		ShowTimestamp   bool
		TimestampFormat string
	}

	tests := []struct {
		name           string
		fields         fields
		message        string
		expectedOutput string
	}{
		{"Basic", fields{false, false, ""}, "ABC", "INFO: ABC\n"},
		{"InfoLevel", fields{true, false, ""}, "ABC", "INFO: ABC\n"},
		{"TimeStamp", fields{true, true, dateFormatString}, "ABC", "INFO: " + currTime + " - ABC\n"},
		{"TimeStamp2", fields{true, true, dateFormatString}, "ABC", "INFO: " + currTime + " - ABC\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = SetLevel("Info")
			f := &VgsTextFormat{
				ShowInfoLevel:   tt.fields.ShowInfoLevel,
				ShowTimestamp:   tt.fields.ShowTimestamp,
				TimestampFormat: tt.fields.TimestampFormat,
			}
			logrus.SetFormatter(f)
			out := CaptureOutput(func() { Logger().Info(tt.message) })
			assert.Equal(t, tt.expectedOutput, out)
		})
	}
	defaultFormat := VgsTextFormat{
		ShowInfoLevel:   false,
		ShowTimestamp:   false,
		TimestampFormat: "",
	}
	lager := Logger()
	lager.Buffer = nil
	retBytes, _ := defaultFormat.Format(lager)
	logs := CaptureOutput(func() {
		lager.Message = "asdf"
		lager.Printf("asfd")
	})
	assert.NotNil(t, retBytes)
	assert.NotNil(t, logs)
}

func TestInitializeLogger(t *testing.T) {
	tests := []struct {
		name      string
		logFormat string
	}{
		{"Text", "text"},
		{"Json", "json"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = os.Setenv("VGS_LOG_FORMAT", test.logFormat)
			logger = nil
			err := initializeLogger()
			assert.NoError(t, err)
		})
	}

	_ = os.Unsetenv("VGS_LOG_FORMAT")

}

func Test_setFormatter(t *testing.T) {
	type args struct {
		layout FormatLayoutType
	}
	tests := []struct {
		name string
		args args
	}{
		{"Text", args{layout: FormatLayoutType("text")}},
		{"Json", args{layout: FormatLayoutType("json")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				setFormatter(tt.args.layout)
			})

		})
	}
}
