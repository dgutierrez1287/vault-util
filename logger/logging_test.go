package logger

import (
  "fmt"
  "github.com/hashicorp/go-hclog"
)

// Mocking the hclog.logger 
type MockLogger struct {
  Level hclog.Level
  Color hclog.ColorOption
}

func (m *MockLogger) Log(level hclog.Level, msg string, args ...interface{}) {
	// No-op for testing, you can expand this to capture logs if needed
}

func (m *MockLogger) IsEnabledFor(level hclog.Level) bool {
	return level >= m.Level
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Debug) {
		fmt.Printf(msg, args...)
	}
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Info) {
		fmt.Printf(msg, args...)
	}
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Warn) {
		fmt.Printf(msg, args...)
	}
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Error) {
		fmt.Printf(msg, args...)
	}
}

func (m *MockLogger) Trace(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Trace) {
		fmt.Printf(msg, args...)
	}
}


