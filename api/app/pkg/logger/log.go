package log

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

const (
	RequestIDField    = "_request_id"
	DebugField        = "_debug"
	VersionField      = "_version"
	DefaultCallerSkip = 2
)

type Logger interface {
	Info(msg string)
	Infof(msg string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Debug(msg string)
	Debugf(msg string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Error(msg string)
	Errorf(msg string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panic(msg string)
	ErrWithError(ctx context.Context, err error, msg string)
	ErrWithErrorf(ctx context.Context, err error, msg string, args ...interface{})
	ErrWithErrorw(ctx context.Context, err error, msg string, keysAndValues ...interface{})
	// LogGRPC integration for grpc unary middleware
	LogGRPC(ctx context.Context, lvl logging.Level, msg string, fields ...any)
}

var (
	glog = Default()
)

func init() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP)

	go func() {
		for {
			<-signals
			glog.Info("Log Rotate signal received")
		}
	}()
}

type Config struct {
	LogLevel         string
	ContextLogFields []string `mapstructure:"context_log_fields"`
	CallerSkip       int

	SentryDSN               string `mapstructure:"sentry_dsn"`
	SentryEnableBreadcrumbs bool
	SentryMaxBreadcrumbs    int

	TgChatID       int64
	TgToken        string
	TgMsgParseMode string
	TgAppName      string
}

type Log struct {
	logger    *zap.SugaredLogger
	Config    Config
	loggerStd *zap.Logger
	debug     bool
}

type Fld map[string]any
type SentryFld map[string]string

func Default() *Log {
	cfgDefault := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfgDefault.Build()
	if err != nil {
		panic(err)
	}

	logger = zap.New(logger.Core(), zap.AddCaller(), zap.AddCallerSkip(DefaultCallerSkip))

	return &Log{
		logger:    logger.Sugar(),
		loggerStd: logger,
		Config: Config{
			ContextLogFields: []string{RequestIDField},
		},
	}
}

func New(cfg Config) *Log {
	l := Default()
	l.Config = cfg
	l.Config.ContextLogFields = addStr(l.Config.ContextLogFields, RequestIDField)

	l.logger = zap.New(l.logger.Desugar().Core(), zap.AddCaller(), zap.AddCallerSkip(cfg.CallerSkip)).Sugar()
	levelLog := l.logger.Level()
	err := levelLog.Set(cfg.LogLevel)
	if err != nil {
		return nil
	}

	return l
}
func (l *Log) GetZapLogger() *zap.Logger {
	return l.loggerStd
}

func (l *Log) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Log) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l *Log) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *Log) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Log) Errorf(msg string, args ...interface{}) {
	l.logger.Errorf(msg, args...)
}

func (l *Log) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *Log) Debug(msg string) {
	if l.debug {
		l.logger.Debug(msg)
	}
}

func (l *Log) Debugf(msg string, args ...interface{}) {
	if l.debug {
		l.logger.Debugf(msg, args...)
	}
}

func (l *Log) Debugw(msg string, keysAndValues ...interface{}) {
	if l.debug {
		l.logger.Debugw(msg, keysAndValues)
	}
}

func (l *Log) With(fld Fld) *Log {
	fields := make([]interface{}, 0)
	for k, v := range fld {
		fields = append(fields, k, v)
	}
	return l.copyWithEntry(*l.logger.With(fields...))
}

func (l *Log) WithField(key string, val interface{}) *Log {
	return l.copyWithEntry(*l.logger).With(Fld{key: val})
}

func (l *Log) WithErr(err error) *Log {
	fields := Fld{}
	if e, ok := err.(errWithFields); ok {
		fields["error"] = e.Origin()
		for k, v := range e.Fields() {
			fields[k] = v
		}
	} else {
		fields["error"] = err
	}

	return l.copyWithEntry(*l.logger).With(fields)
}

func (l *Log) ErrWithError(ctx context.Context, err error, msg string) {
	l.WithCtx(ctx).WithErr(err).Error(msg)
}

func (l *Log) ErrWithErrorf(ctx context.Context, err error, msg string, args ...interface{}) {
	l.WithCtx(ctx).WithErr(err).Errorf(msg, args...)
}

func (l *Log) ErrWithErrorw(ctx context.Context, err error, msg string, keysAndValues ...interface{}) {
	l.WithCtx(ctx).WithErr(err).Errorw(msg, keysAndValues...)
}

func (l *Log) WithCtx(ctx context.Context) *Log {
	fields := Fld{}
	for _, key := range l.Config.ContextLogFields {
		v := ctx.Value(key)
		if v != nil {
			fields[key] = v
		}
	}

	copied := l.copyWithEntry(*l.logger).With(fields)
	if ctx.Value(DebugField) != nil {
		copied.debug = true
	}

	return copied
}
func (l *Log) Panic(msg string) {
	l.logger.Panic(msg)
}

func (l *Log) LogGRPC(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
	switch lvl {
	case logging.LevelDebug:
		l.logger.Debugw(msg, fields...)
	case logging.LevelInfo:
		l.logger.Infow(msg, fields...)
	case logging.LevelWarn:
		l.logger.Warnw(msg, fields...)
	case logging.LevelError:
		l.logger.Errorw(msg, fields...)
	default:
		panic(fmt.Sprintf("unknown level %v", lvl))
	}
}

func (l *Log) copyWithEntry(entry zap.SugaredLogger) *Log {
	return &Log{
		logger: &entry,
		Config: l.Config,
	}
}

// LogPanic логирует пойманную панику
func LogPanic(recovered interface{}) { // nolint: revive
	if recovered == nil {
		return
	}
	glog.Errorf("Panic recovered: %s %s", recovered, string(debug.Stack()))
}

func (l *Log) LogPanic(recovered interface{}) {
	if recovered == nil {
		return
	}
	l.Errorf("Panic recovered: %s %s", recovered, string(debug.Stack()))
}

func (l *Log) Log(ctx context.Context, level int, msg string, fields ...interface{}) {
	l.WithCtx(ctx).logger.Logf(zapcore.Level(level), msg, fields...)
}

func addStr(ss []string, s string) []string {
	for i := range ss {
		if ss[i] == s {
			return ss
		}
	}

	return append(ss, s)
}
