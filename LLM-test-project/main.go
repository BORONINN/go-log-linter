package main

import (
	"context"
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	// Инициализация логгеров
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	zapSugar := zapLogger.Sugar()

	ctx := context.Background()

	// ========================================================================
	// ТЕСТ 1: ПРАВИЛО lowercase (строчная буква)
	// ========================================================================
	println("\n=== ТЕСТ 1: lowercase rule ===")

	// ДОЛЖНЫ ВЫЗВАТЬ ОШИБКИ (начинаются с заглавной)
	slog.Info("Starting server")             // ОШИБКА
	slog.Error("Failed to connect")          // ОШИБКА
	slog.Warn("Warning issued")              // ОШИБКА
	slog.Debug("Debug message")              // ОШИБКА
	slog.InfoContext(ctx, "Request started") // ОШИБКА

	// ДОЛЖНЫ ВЫЗВАТЬ ОШИБКИ (zap)
	zapLogger.Info("Server started")     // ОШИБКА
	zapLogger.Error("Database error")    // ОШИБКА
	zapSugar.Infow("User logged in")     // ОШИБКА
	zapSugar.Errorw("Operation failed")  // ОШИБКА
	zapSugar.Infof("Processing request") // ОШИБКА
	zapSugar.Errorf("Invalid input")     // ОШИБКА

	// НЕ ДОЛЖНЫ ВЫЗЫВАТЬ ОШИБОК (начинаются со строчной)
	slog.Info("starting server")
	slog.Error("failed to connect")
	slog.Warn("warning issued")
	slog.Debug("debug message")
	slog.InfoContext(ctx, "request started")

	zapLogger.Info("server started")
	zapLogger.Error("database error")
	zapSugar.Infow("user logged in")
	zapSugar.Errorw("operation failed")
	zapSugar.Infof("processing request")
	zapSugar.Errorf("invalid input")

	// ========================================================================
	// ТЕСТ 2: ПРАВИЛО english (только английский)
	// ========================================================================
	println("\n=== ТЕСТ 2: english rule ===")

	// ДОЛЖНЫ ВЫЗВАТЬ ОШИБКИ (не английские буквы)
	slog.Info("запуск сервера")      // ОШИБКА (русский)
	slog.Error("ошибка подключения") // ОШИБКА (русский)
	slog.Info("服务器启动")               // ОШИБКА (китайский)
	slog.Warn("サーバー起動")              // ОШИБКА (японский)

	zapLogger.Info("запуск сервера")     // ОШИБКА (русский)
	zapSugar.Infow("ошибка базы данных") // ОШИБКА (русский)

	// НЕ ДОЛЖНЫ ВЫЗЫВАТЬ ОШИБОК (английский, включая диакритику)
	slog.Info("starting server")
	slog.Info("café is ready")  // OK (é - латинская)
	slog.Info("über service")   // OK (ü - латинская)
	slog.Info("façade pattern") // OK (ç - латинская)

	zapLogger.Info("server started")
	zapSugar.Infow("request completed")

	// ========================================================================
	// ТЕСТ 3: ПРАВИЛО special-chars (спецсимволы и эмодзи)
	// ========================================================================
	println("\n=== ТЕСТ 3: special-chars rule ===")

	// ДОЛЖНЫ ВЫЗВАТЬ ОШИБКИ (содержат спецсимволы)
	slog.Info("server started!")         // ОШИБКА (!)
	slog.Error("connection failed!!!")   // ОШИБКА (!!!)
	slog.Warn("something went wrong...") // ОШИБКА (...)
	slog.Info("deployed 🚀")              // ОШИБКА (эмодзи)
	slog.Debug("hot 🔥")                  // ОШИБКА (эмодзи)

	zapLogger.Info("process finished!") // ОШИБКА (!)
	zapSugar.Infow("task completed...") // ОШИБКА (...)
	zapSugar.Infof("result: %d!!!", 42) // ОШИБКА (!!!)

	// НЕ ДОЛЖНЫ ВЫЗЫВАТЬ ОШИБОК (без спецсимволов)
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("something went wrong")
	slog.Info("version 1.2.3") // OK (точки в версии)
	slog.Info("key: value")    // OK (двоеточие)
	slog.Info("a,b,c")         // OK (запятые)

	zapLogger.Info("process finished")
	zapSugar.Infow("task completed")

	// ========================================================================
	// ТЕСТ 4: ПРАВИЛО sensitive (чувствительные данные)
	// ========================================================================
	println("\n=== ТЕСТ 4: sensitive rule ===")

	// Переменные с чувствительными данными
	password := "secret123"
	apiKey := "abc-123-def"
	token := "jwt-token-here"
	ssn := "123-45-6789"

	// ДОЛЖНЫ ВЫЗВАТЬ ОШИБКИ (ключевые слова в тексте)
	slog.Info("user password: admin123") // ОШИБКА (password)
	slog.Info("api_key = abc123")        // ОШИБКА (api_key)
	slog.Debug("token: xyz789")          // ОШИБКА (token)
	slog.Error("invalid secret")         // ОШИБКА (secret)

	// ДОЛЖНЫ ВЫЗВАТЬ ОШИБКИ (конкатенация с переменными)
	slog.Info("password: " + password) // ОШИБКА (password)
	slog.Info("api_key=" + apiKey)     // ОШИБКА (api_key)
	slog.Debug("token=" + token)       // ОШИБКА (token)
	slog.Info("ssn: " + ssn)           // ОШИБКА (ssn - если добавили флаг)

	// ДОЛЖНЫ ВЫЗВАТЬ ОШИБКИ (zap)
	zapLogger.Info("user password entered")   // ОШИБКА (password)
	zapSugar.Infow("token validation failed") // ОШИБКА (token)
	zapSugar.Infof("api_key: %s", apiKey)     // ОШИБКА (api_key)

	// НЕ ДОЛЖНЫ ВЫЗЫВАТЬ ОШИБОК (безопасные сообщения)
	slog.Info("user authenticated successfully")
	slog.Debug("request completed")
	slog.Error("database connection timeout")
	slog.Info("cache miss", "key", "user:123")

	zapLogger.Info("user logged in", zap.String("user_id", "12345"))
	zapSugar.Infow("request processed", "status", 200, "duration_ms", 42)

	// ========================================================================
	// ТЕСТ 5: КРАЕВЫЕ СЛУЧАИ
	// ========================================================================
	println("\n=== ТЕСТ 5: edge cases ===")

	// Смешанные случаи (должна сработать lowercase)
	slog.Info("SERVER STARTED! 🚀") // ОШИБКА (lowercase)

	// Пустые сообщения (не должны вызывать ошибок)
	slog.Info("") // OK

	// Сообщения, начинающиеся с цифр (не должны вызывать lowercase)
	slog.Info("123 items processed")     // OK
	slog.Info("2 factor authentication") // OK

	// Сообщения с безопасными спецсимволами
	slog.Info("file: /var/log/app.log")  // OK (слеш)
	slog.Info("email: user@example.com") // OK (@)
	slog.Info("path: C:\\Program Files") // OK (обратный слеш)

	// ========================================================================
	// ТЕСТ 6: НЕ-ЛИТЕРАЛЫ (должны пропускаться)
	// ========================================================================
	println("\n=== ТЕСТ 6: non-literals (should be skipped) ===")

	msgVar := "dynamic message"
	getMsg := func() string { return "hello" }

	slog.Info(msgVar)              // SKIP (переменная)
	slog.Info(getMsg())            // SKIP (вызов функции)
	slog.Info("prefix: " + msgVar) // SKIP (конкатенация с переменной)
}

// Дополнительная функция для тестирования разных уровней логирования
func testLevels(logger *zap.Logger, sugar *zap.SugaredLogger) {
	// Все эти вызовы должны проверяться одинаково
	logger.Debug("Debug message")   // ОШИБКА (если uppercase)
	logger.Warn("Warn message")     // ОШИБКА (если uppercase)
	logger.Error("Error message")   // ОШИБКА (если uppercase)
	logger.Fatal("Fatal message")   // ОШИБКА (если uppercase)
	logger.Panic("Panic message")   // ОШИБКА (если uppercase)
	logger.DPanic("DPanic message") // ОШИБКА (если uppercase)

	sugar.Debugw("Debugw message") // ОШИБКА (если uppercase)
	sugar.Warnw("Warnw message")   // ОШИБКА (если uppercase)
	sugar.Fatalw("Fatalw message") // ОШИБКА (если uppercase)
	sugar.Panicw("Panicw message") // ОШИБКА (если uppercase)
}
