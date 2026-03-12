// Package zap is a minimal stub of go.uber.org/zap for testing.
package zap

// Field is a stub for zap.Field.
type Field struct{}

// String creates a string Field stub.
func String(_, _ string) Field { return Field{} }

// Int creates an int Field stub.
func Int(_ string, _ int) Field { return Field{} }

// Error creates an error Field stub.
func Error(_ error) Field { return Field{} }

// Logger is a stub for *zap.Logger.
type Logger struct{}

// NewNop creates a no-op Logger stub.
func NewNop() *Logger { return &Logger{} }

// Sugar returns a SugaredLogger stub.
func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

// Logger methods — structured logging with typed fields.

// Info is a stub.
func (l *Logger) Info(_ string, _ ...Field) {}

// Error is a stub.
func (l *Logger) Error(_ string, _ ...Field) {}

// Debug is a stub.
func (l *Logger) Debug(_ string, _ ...Field) {}

// Warn is a stub.
func (l *Logger) Warn(_ string, _ ...Field) {}

// Fatal is a stub.
func (l *Logger) Fatal(_ string, _ ...Field) {}

// Panic is a stub.
func (l *Logger) Panic(_ string, _ ...Field) {}

// DPanic is a stub.
func (l *Logger) DPanic(_ string, _ ...Field) {}

// SugaredLogger is a stub for *zap.SugaredLogger.
type SugaredLogger struct{}

// SugaredLogger methods — printf-style and key-value-style logging.

// Infow is a stub.
func (s *SugaredLogger) Infow(_ string, _ ...interface{}) {}

// Errorw is a stub.
func (s *SugaredLogger) Errorw(_ string, _ ...interface{}) {}

// Debugw is a stub.
func (s *SugaredLogger) Debugw(_ string, _ ...interface{}) {}

// Warnw is a stub.
func (s *SugaredLogger) Warnw(_ string, _ ...interface{}) {}

// Fatalw is a stub.
func (s *SugaredLogger) Fatalw(_ string, _ ...interface{}) {}

// Panicw is a stub.
func (s *SugaredLogger) Panicw(_ string, _ ...interface{}) {}

// Infof is a stub.
func (s *SugaredLogger) Infof(_ string, _ ...interface{}) {}

// Errorf is a stub.
func (s *SugaredLogger) Errorf(_ string, _ ...interface{}) {}

// Debugf is a stub.
func (s *SugaredLogger) Debugf(_ string, _ ...interface{}) {}

// Warnf is a stub.
func (s *SugaredLogger) Warnf(_ string, _ ...interface{}) {}

// Fatalf is a stub.
func (s *SugaredLogger) Fatalf(_ string, _ ...interface{}) {}

// Panicf is a stub.
func (s *SugaredLogger) Panicf(_ string, _ ...interface{}) {}
