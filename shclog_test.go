package shclog

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShclog(t *testing.T) {
	t.Run("logging", func(t *testing.T) {
		output := &bytes.Buffer{}
		slogger := slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: SlogLevelTrace,
		}))
		l := New(slogger)

		l.Log(hclog.Off, "message -1")
		l.Log(hclog.Debug, "message 0", "key0", "value0")
		l.Trace("message 1", "key1", "value1")
		l.Debug("message 2", "key2", "value2")
		l.Info("message 3", "key3", "value3")
		l.Warn("message 4", "key4", "value4")
		l.Error("message 5", "key5", "value5")

		o := output.String()

		require.NotContains(t, o, `message -1`)
		require.Contains(t, o, `level=DEBUG msg="message 0" key0=value0`)
		require.Contains(t, o, `level=DEBUG-1 msg="message 1" key1=value1`)
		require.Contains(t, o, `level=DEBUG msg="message 2" key2=value2`)
		require.Contains(t, o, `level=INFO msg="message 3" key3=value3`)
		require.Contains(t, o, `level=WARN msg="message 4" key4=value4`)
		require.Contains(t, o, `level=ERROR msg="message 5" key5=value5`)
	})

	t.Run("naming", func(t *testing.T) {
		output := &bytes.Buffer{}
		slogger := slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		l := New(slogger)

		l = l.Named("name0")
		l = l.Named("name1")

		l.Info("message 1")

		l = l.ResetNamed("name2")

		l.Info("message 2")

		l = l.ResetNamed("")

		l.Info("message 3")

		o := output.String()

		require.Contains(t, o, `level=INFO msg="name0: name1: message 1"`)
		require.Contains(t, o, `level=INFO msg="name2: message 2"`)
		require.Contains(t, o, `level=INFO msg="message 3"`)
	})

	t.Run("fielding", func(t *testing.T) {
		output := &bytes.Buffer{}
		slogger := slog.New(slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		l := New(slogger)
		l = l.With("key1", "value1")
		l = l.With("key2", "value2")

		l.Info("message 1")

		args := l.ImpliedArgs()

		o := output.String()

		require.Contains(t, o, `level=INFO msg="message 1" key1=value1 key2=value2`)
		require.Contains(t, args, "key1", "value1", "key2", "value2")
	})

	t.Run("getting log level", func(t *testing.T) {
		slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		l := New(slogger)

		isTrace := l.IsTrace()
		isDebug := l.IsDebug()
		isInfo := l.IsInfo()
		isWarn := l.IsWarn()
		isError := l.IsError()

		assert.Equal(t, false, isTrace)
		assert.Equal(t, false, isDebug)
		assert.Equal(t, true, isInfo)
		assert.Equal(t, false, isWarn)
		assert.Equal(t, false, isError)
	})

	t.Run("std", func(t *testing.T) {
		slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
		l := New(slogger)
		stdLog := l.StandardLogger(&hclog.StandardLoggerOptions{})
		assert.NotNil(t, stdLog)
		stdLog.Println("message 1")

		w := l.StandardWriter(&hclog.StandardLoggerOptions{})
		assert.NotNil(t, w)
		_, err := w.Write([]byte("message 2"))
		assert.Nil(t, err)
	})
}
