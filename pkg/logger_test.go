package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInit(t *testing.T) {
	Init()
	if Log == nil {
		t.Error("Log should not be nil after Init()")
	}
}

func TestLoggerMethods(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a custom logger for testing
	testEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	testCore := zapcore.NewCore(testEncoder, zapcore.AddSync(&buf), zapcore.InfoLevel)
	testLogger := &Logger{zap.New(testCore)}

	tests := []struct {
		name     string
		logFunc  func(string, ...zap.Field)
		message  string
		fields   []zap.Field
		expected string
	}{
		{"Info", testLogger.Info, "info message", []zap.Field{zap.String("key", "value")}, "info"},
		{"Error", testLogger.Error, "error message", []zap.Field{zap.Int("code", 500)}, "error"},
		{"With", func(msg string, fields ...zap.Field) {
			testLogger.With(zap.String("with", "field")).Info(msg, fields...)
		}, "with message", nil, "info"}, // Changed expected level to "info"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message, tt.fields...)

			var logEntry map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
				t.Fatalf("Failed to unmarshal log entry: %v", err)
			}

			if logEntry["level"] != tt.expected {
				t.Errorf("Expected log level %s, got %s", tt.expected, logEntry["level"])
			}
			if logEntry["msg"] != tt.message {
				t.Errorf("Expected message %s, got %s", tt.message, logEntry["msg"])
			}

			for _, field := range tt.fields {
				if value, ok := logEntry[field.Key]; !ok {
					t.Errorf("Expected field %s not found in log entry", field.Key)
				} else {
					// Compare values based on the field type
					switch field.Type {
					case zapcore.StringType:
						if value != field.String {
							t.Errorf("Expected field %s with value %v, got %v", field.Key, field.String, value)
						}
					case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
						if int64(value.(float64)) != field.Integer {
							t.Errorf("Expected field %s with value %v, got %v", field.Key, field.Integer, value)
						}
					default:
						// For other types, use string comparison
						if fmt.Sprintf("%v", value) != fmt.Sprintf("%v", field.Interface) {
							t.Errorf("Expected field %s with value %v, got %v", field.Key, field.Interface, value)
						}
					}
				}
			}

			if tt.name == "With" {
				if value, ok := logEntry["with"]; !ok || value != "field" {
					t.Errorf("Expected 'with' field with value 'field', got %v", value)
				}
			}
		})
	}
}

func TestGlobalLoggerMethods(t *testing.T) {
	// Initialize the global logger
	Init()

	// Redirect global logger output to a buffer
	var buf bytes.Buffer
	Log.Logger = zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.AddSync(&buf),
			zapcore.InfoLevel,
		),
	)

	tests := []struct {
		name    string
		logFunc func(string, ...zap.Field)
		message string
		fields  []zap.Field
	}{
		{"Info", Info, "global info", []zap.Field{zap.String("global", "value")}},
		{"Error", Error, "global error", []zap.Field{zap.Int("global_code", 500)}},
		{"With", func(msg string, fields ...zap.Field) {
			With(zap.String("global_with", "field")).Info(msg, fields...)
		}, "global with", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.logFunc(tt.message, tt.fields...)

			logOutput := buf.String()
			if !strings.Contains(logOutput, tt.message) {
				t.Errorf("Log output doesn't contain expected message: %s", tt.message)
			}

			for _, field := range tt.fields {
				if !strings.Contains(logOutput, field.Key) || !strings.Contains(logOutput, field.String) {
					t.Errorf("Log output doesn't contain expected field: %s", field.Key)
				}
			}

			if tt.name == "With" && !strings.Contains(logOutput, "global_with") {
				t.Errorf("Log output doesn't contain 'global_with' field")
			}
		})
	}
}
