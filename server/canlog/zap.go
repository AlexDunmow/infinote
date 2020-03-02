package canlog

import "go.uber.org/zap"

// ZapLogger holds context for single line logs
type ZapLogger struct {
	*zap.Logger
	DetectErr bool
}

// Log the contexted logger to single line
func (l *ZapLogger) Log(msg string, fields map[string]interface{}) {
	if l.DetectErr {
		_, ok := fields["err"]
		if ok {
			l.Error(msg, zap.Reflect("fields", fields))
			return
		}
	}

	l.Info(msg, zap.Reflect("fields", fields))
}
