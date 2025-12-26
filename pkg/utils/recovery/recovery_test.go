package recovery

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// setupTestLogger 创建测试用的logger / Create test logger
func setupTestLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	return logger, logs
}

func TestNewRecoveryHandler(t *testing.T) {
	logger, _ := setupTestLogger()
	handler := NewRecoveryHandler(logger)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
}

func TestRecover_NoPanic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	err := handler.Recover(func() error {
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 0, logs.Len(), "should not log when no panic")
}

func TestRecover_WithError(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	expectedErr := errors.New("test error")
	err := handler.Recover(func() error {
		return expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 0, logs.Len(), "should not log when returning error normally")
}

func TestRecover_PanicWithString(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	err := handler.Recover(func() error {
		panic("test panic")
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic: test panic")
	assert.Equal(t, 1, logs.Len(), "should log panic")

	logEntry := logs.All()[0]
	assert.Equal(t, "panic recovered", logEntry.Message)
	assert.Equal(t, "test panic", logEntry.ContextMap()["panic"])
}

func TestRecover_PanicWithError(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	panicErr := errors.New("panic error")
	err := handler.Recover(func() error {
		panic(panicErr)
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic:")
	assert.ErrorIs(t, err, panicErr)
	assert.Equal(t, 1, logs.Len(), "should log panic")
}

func TestRecover_PanicWithInt(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	err := handler.Recover(func() error {
		panic(42)
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic: 42")
	assert.Equal(t, 1, logs.Len(), "should log panic")
}

func TestRecoverWithValue_NoPanic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	result, err := handler.RecoverWithValue(func() (interface{}, error) {
		return "success", nil
	})

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
	assert.Equal(t, 0, logs.Len(), "should not log when no panic")
}

func TestRecoverWithValue_WithError(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	expectedErr := errors.New("test error")
	result, err := handler.RecoverWithValue(func() (interface{}, error) {
		return nil, expectedErr
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, logs.Len(), "should not log when returning error normally")
}

func TestRecoverWithValue_Panic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	result, err := handler.RecoverWithValue(func() (interface{}, error) {
		panic("test panic")
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic: test panic")
	assert.Nil(t, result)
	assert.Equal(t, 1, logs.Len(), "should log panic")
}

func TestSafeGo_NoPanic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	done := make(chan bool)
	handler.SafeGo(func() {
		done <- true
	})

	select {
	case <-done:
		// Success
	case <-time.After(time.Second):
		t.Fatal("goroutine did not complete")
	}

	assert.Equal(t, 0, logs.Len(), "should not log when no panic")
}

func TestSafeGo_Panic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	handler.SafeGo(func() {
		panic("goroutine panic")
	})

	// 等待一小段时间让goroutine执行 / Wait a bit for goroutine to execute
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, logs.Len(), "should log panic in goroutine")
	logEntry := logs.All()[0]
	assert.Equal(t, "goroutine panic recovered", logEntry.Message)
}

func TestWrapHandler_NoPanic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	originalHandler := func() error {
		return nil
	}

	wrappedHandler := handler.WrapHandler(originalHandler)
	err := wrappedHandler()

	assert.NoError(t, err)
	assert.Equal(t, 0, logs.Len(), "should not log when no panic")
}

func TestWrapHandler_Panic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	originalHandler := func() error {
		panic("wrapped panic")
	}

	wrappedHandler := handler.WrapHandler(originalHandler)
	err := wrappedHandler()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic: wrapped panic")
	assert.Equal(t, 1, logs.Len(), "should log panic")
}

func TestWrapHandlerWithValue_NoPanic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	originalHandler := func() (interface{}, error) {
		return "result", nil
	}

	wrappedHandler := handler.WrapHandlerWithValue(originalHandler)
	result, err := wrappedHandler()

	assert.NoError(t, err)
	assert.Equal(t, "result", result)
	assert.Equal(t, 0, logs.Len(), "should not log when no panic")
}

func TestWrapHandlerWithValue_Panic(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	originalHandler := func() (interface{}, error) {
		panic("wrapped panic with value")
	}

	wrappedHandler := handler.WrapHandlerWithValue(originalHandler)
	result, err := wrappedHandler()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "panic: wrapped panic with value")
	assert.Nil(t, result)
	assert.Equal(t, 1, logs.Len(), "should log panic")
}

func TestPanicToError_DifferentTypes(t *testing.T) {
	logger, _ := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	tests := []struct {
		name        string
		panicValue  interface{}
		expectedMsg string
	}{
		{
			name:        "string panic",
			panicValue:  "string error",
			expectedMsg: "panic: string error",
		},
		{
			name:        "error panic",
			panicValue:  errors.New("error type"),
			expectedMsg: "panic: error type",
		},
		{
			name:        "int panic",
			panicValue:  123,
			expectedMsg: "panic: 123",
		},
		{
			name:        "nil panic",
			panicValue:  nil,
			expectedMsg: "panic: <nil>",
		},
		{
			name:        "struct panic",
			panicValue:  struct{ msg string }{msg: "test"},
			expectedMsg: "panic: {test}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.panicToError(tt.panicValue)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedMsg)
		})
	}
}

func TestRecover_StackTraceLogged(t *testing.T) {
	logger, logs := setupTestLogger()
	handler := NewRecoveryHandler(logger)

	err := handler.Recover(func() error {
		panic("test stack trace")
	})

	require.Error(t, err)
	assert.Equal(t, 1, logs.Len(), "should log panic")

	logEntry := logs.All()[0]
	stackTrace, ok := logEntry.ContextMap()["stack"].(string)
	require.True(t, ok, "stack trace should be logged")
	assert.NotEmpty(t, stackTrace, "stack trace should not be empty")
	assert.Contains(t, stackTrace, "recovery_test.go", "stack trace should contain file name")
}
