// Package zapcases tests the logcheck analyzer with go.uber.org/zap calls.
package zapcases

import "go.uber.org/zap"

func zapLoggerExamples() {
	logger := zap.NewNop()

	// Lowercase rule.
	logger.Info("Starting server") // want `log message should start with a lowercase letter`
	logger.Error("Error occurred") // want `log message should start with a lowercase letter`
	logger.Info("starting server") // OK
	logger.Error("error occurred") // OK

	// English rule.
	logger.Info("запуск сервера")  // want `log message should be in English`
	logger.Info("starting server") // OK

	// Special chars rule.
	logger.Info("done! 🚀")    // want `log message should not contain special character`
	logger.Warn("loading...") // want `log message should not contain`
	logger.Info("done")       // OK
}

func zapSugaredExamples() {
	logger := zap.NewNop()
	sugar := logger.Sugar()

	// Lowercase rule.
	sugar.Infow("Starting request") // want `log message should start with a lowercase letter`
	sugar.Errorw("Database error")  // want `log message should start with a lowercase letter`
	sugar.Infow("starting request") // OK
	sugar.Errorw("database error")  // OK

	// English rule.
	sugar.Infow("ошибка запроса") // want `log message should be in English`
	sugar.Infow("request error")  // OK

	// Special chars.
	sugar.Infow("success!") // want `log message should not contain special character`
	sugar.Infow("success")  // OK

	// Sensitive data.
	token := "abc"
	sugar.Infow("token: " + token)   // want `log message may contain sensitive data`
	sugar.Infow("request completed") // OK
}

func zapFormatExamples() {
	logger := zap.NewNop()
	sugar := logger.Sugar()

	sugar.Infof("Starting service")  // want `log message should start with a lowercase letter`
	sugar.Errorf("Error in handler") // want `log message should start with a lowercase letter`
	sugar.Infof("starting service")  // OK
}

func zapFieldExamples() {
	logger := zap.NewNop()

	// Fields should not affect message checking.
	logger.Info("Starting with fields", zap.String("key", "val")) // want `log message should start with a lowercase letter`
	logger.Info("starting with fields", zap.String("key", "val")) // OK
	logger.Info("starting with int", zap.Int("count", 42))        // OK
}
