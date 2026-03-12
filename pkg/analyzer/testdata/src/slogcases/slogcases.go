// Package slogcases tests the logcheck analyzer with log/slog calls.
package slogcases

import "log/slog"

func lowercaseExamples() {
	slog.Info("Starting server on port 8080") // want `log message should start with a lowercase letter`
	slog.Error("Failed to connect")           // want `log message should start with a lowercase letter`
	slog.Debug("Debug info here")             // want `log message should start with a lowercase letter`
	slog.Warn("Warning issued")               // want `log message should start with a lowercase letter`

	slog.Info("starting server on port 8080") // OK
	slog.Error("failed to connect")           // OK
	slog.Debug("debug info here")             // OK
	slog.Warn("warning issued")               // OK
	slog.Info("123 items processed")          // OK — starts with digit
}

func englishExamples() {
	slog.Info("запуск сервера")             // want `log message should be in English`
	slog.Error("ошибка подключения к базе") // want `log message should be in English`
	slog.Warn("内部エラー")                      // want `log message should be in English`

	slog.Info("starting server") // OK
	slog.Info("café opened")     // OK — latin diacritics are fine
}

func specialCharsExamples() {
	slog.Info("server started!")         // want `log message should not contain special character`
	slog.Error("connection failed!!!")   // want `log message should not contain special character`
	slog.Warn("something went wrong...") // want `log message should not contain`
	slog.Info("deployed 🚀")              // want `log message should not contain emoji`

	slog.Info("server started") // OK
	slog.Info("version 1.2.3")  // OK — single dots are fine
	slog.Info("key: value")     // OK — colons are fine
}

func sensitiveExamples() {
	password := "secret123"
	apiKey := "key-abc-123"
	token := "tok-xyz"

	slog.Info("user password: " + password) // want `log message may contain sensitive data`
	slog.Debug("api_key=" + apiKey)         // want `log message may contain sensitive data`
	slog.Info("token: " + token)            // want `log message may contain sensitive data`

	slog.Info("user authenticated successfully") // OK
	slog.Debug("api request completed")          // OK
	slog.Info("operation finished")              // OK
}

func loggerMethodExamples() {
	logger := slog.New(nil)

	logger.Info("Starting background job") // want `log message should start with a lowercase letter`
	logger.Error("Failed to save")         // want `log message should start with a lowercase letter`

	logger.Info("starting background job") // OK
	logger.Error("failed to save")         // OK
}
