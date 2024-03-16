package shclog

import (
	"context"
	"io"
	"log"
	"log/slog"

	"github.com/hashicorp/go-hclog"
)

// Slog does not have built in trace log level
const SlogLevelTrace = slog.LevelDebug - 4

type Shclog struct {
	l *slog.Logger

	name        string
	level       hclog.Level
	impliedArgs []interface{}
}

func New(l *slog.Logger) hclog.Logger {
	return &Shclog{l: l, level: getSlogLevel(l)}
}

func (l *Shclog) Log(level hclog.Level, msg string, args ...interface{}) {
	switch level {
	case hclog.Trace:
		l.Trace(msg, args...)
	case hclog.Debug:
		l.Debug(msg, args...)
	case hclog.Info:
		l.Info(msg, args...)
	case hclog.Warn:
		l.Warn(msg, args...)
	case hclog.Error:
		l.Error(msg, args...)
	case hclog.Off:
	default:
		l.Info(msg, args...)
	}
}

func (l *Shclog) Trace(msg string, args ...interface{}) {
	l.l.Log(context.Background(), SlogLevelTrace, msg, args...)
}

func (l *Shclog) Debug(msg string, args ...interface{}) {
	l.l.Debug(l.name+msg, args...)
}

func (l *Shclog) Info(msg string, args ...interface{}) {
	l.l.Info(l.name+msg, args...)
}

func (l *Shclog) Warn(msg string, args ...interface{}) {
	l.l.Warn(l.name+msg, args...)
}

func (l *Shclog) Error(msg string, args ...interface{}) {
	l.l.Error(l.name+msg, args...)
}

func (l *Shclog) IsTrace() bool {
	return l.level == hclog.Trace
}

func (l *Shclog) IsDebug() bool {
	return l.level == hclog.Debug
}

func (l *Shclog) IsInfo() bool {
	return l.level == hclog.Info
}

func (l *Shclog) IsWarn() bool {
	return l.level == hclog.Warn
}

func (l *Shclog) IsError() bool {
	return l.level == hclog.Error
}

func (l *Shclog) ImpliedArgs() []interface{} {
	return l.impliedArgs
}

func (l *Shclog) With(args ...interface{}) hclog.Logger {
	sl := cloneSlog(l.l)
	impliedArgs := append(l.impliedArgs, args...)

	return &Shclog{
		l:           sl.With(impliedArgs...),
		name:        l.name,
		level:       l.level,
		impliedArgs: impliedArgs,
	}
}

func (l *Shclog) Name() string {
	if l.name == "" {
		return ""
	}
	return l.name[:len(l.name)-2]
}

func (l *Shclog) Named(name string) hclog.Logger {
	sl := cloneSlog(l.l)

	return &Shclog{
		l:           sl,
		name:        l.name + name + ": ",
		level:       l.level,
		impliedArgs: l.impliedArgs,
	}
}

func (l *Shclog) ResetNamed(name string) hclog.Logger {
	sl := cloneSlog(l.l)

	return &Shclog{
		l:           sl,
		name:        name + ": ",
		level:       l.level,
		impliedArgs: l.impliedArgs,
	}
}

func (l *Shclog) SetLevel(level hclog.Level) {
	//noop: Can not set level from here. Please set it through slog.Handler options
}

func (l *Shclog) GetLevel() hclog.Level {
	return l.level
}

func (l *Shclog) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	if opts == nil {
		opts = &hclog.StandardLoggerOptions{}
	}
	return log.New(l.StandardWriter(opts), "", 0)
}

func (l *Shclog) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	newLog := cloneShclog(l)

	return &stdlogAdapter{
		log:                      newLog,
		inferLevels:              opts.InferLevels,
		inferLevelsWithTimestamp: opts.InferLevelsWithTimestamp,
		forceLevel:               opts.ForceLevel,
	}
}

func getSlogLevel(l *slog.Logger) hclog.Level {
	h := l.Handler()
	ctx := context.Background()

	if h.Enabled(ctx, SlogLevelTrace) {
		return hclog.Trace
	}
	if h.Enabled(ctx, slog.LevelDebug) {
		return hclog.Debug
	}
	if h.Enabled(ctx, slog.LevelInfo) {
		return hclog.Info
	}
	if h.Enabled(ctx, slog.LevelError) {
		return hclog.Error
	}

	return hclog.Info
}

func cloneSlog(l *slog.Logger) *slog.Logger {
	c := *l
	return &c
}

func cloneShclog(l *Shclog) *Shclog {
	c := cloneSlog(l.l)
	return &Shclog{
		l:           c,
		name:        l.name,
		level:       l.level,
		impliedArgs: l.impliedArgs,
	}
}
