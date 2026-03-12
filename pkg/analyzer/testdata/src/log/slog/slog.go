// Package slog is a minimal stub of log/slog for testing the logcheck analyzer.
//
// Only the signatures matter — the type checker needs to resolve calls like
// slog.Info("msg") to the correct package path and function signature.
// The bodies are empty because the analyzer only inspects the AST and types,
// it never executes the code.
package slog

// Handler is a stub interface.
type Handler interface{}

// Logger is a stub for *slog.Logger.
type Logger struct{}

// New returns a new Logger stub.
func New(_ Handler) *Logger { return &Logger{} }

// With returns a Logger stub.
func (l *Logger) With(_ ...any) *Logger { return l }

// Package-level functions — these are the most common usage pattern.

// Info is a stub.
func Info(_ string, _ ...any) {}

// Error is a stub.
func Error(_ string, _ ...any) {}

// Debug is a stub.
func Debug(_ string, _ ...any) {}

// Warn is a stub.
func Warn(_ string, _ ...any) {}

// InfoContext is a stub.
func InfoContext(_ any, _ string, _ ...any) {}

// ErrorContext is a stub.
func ErrorContext(_ any, _ string, _ ...any) {}

// DebugContext is a stub.
func DebugContext(_ any, _ string, _ ...any) {}

// WarnContext is a stub.
func WarnContext(_ any, _ string, _ ...any) {}

// Methods on *Logger — same names, but called on a receiver.

// Info is a stub method.
func (l *Logger) Info(_ string, _ ...any) {}

// Error is a stub method.
func (l *Logger) Error(_ string, _ ...any) {}

// Debug is a stub method.
func (l *Logger) Debug(_ string, _ ...any) {}

// Warn is a stub method.
func (l *Logger) Warn(_ string, _ ...any) {}
