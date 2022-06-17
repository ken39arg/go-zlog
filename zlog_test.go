package zlog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"testing"

	"github.com/ken39arg/go-zlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestWrappers(t *testing.T) {
	t.Run("safety not initialized", func(t *testing.T) {
		zlog.Infof(context.Background(), "foo %s", "bar")
	})

	logger, logread := testLoggerAndReader()
	defer zlog.SetGlobalLogger(logger.Named("top"))()

	t.Run("context", func(t *testing.T) {
		ctx := zlog.CtxWithLogger(context.Background(), zlog.Logger(context.Background()).Named("child"))

		zlog.Infof(context.Background(), "oh my %s", "god")
		l := logread(t)
		if l["N"] != "top" {
			t.Error("parent context must be use global logger")
		}

		zlog.Infof(ctx, "oh %s", "no")
		l = logread(t)
		if l["N"] != "top.child" {
			t.Error("innner context must be use child logger")
		}
	})

	t.Run("leveld", func(t *testing.T) {
		ctx := context.Background()

		for _, tc := range []struct {
			l string
			f func(ctx context.Context, template string, args ...interface{})
		}{
			{"DEBUG", zlog.Debugf},
			{"INFO", zlog.Infof},
			{"WARN", zlog.Warnf},
			{"ERROR", zlog.Errorf},
			{"PANIC", zlog.Panicf},
		} {
			t.Run(tc.l, func(t *testing.T) {
				n := rand.Intn(1000)
				defer func() {
					if e := recover(); e != nil {
						log.Printf("recover %s", e)
					}
					v := logread(t)
					if e := fmt.Sprintf("%s %d", tc.l, n); e != v["M"] {
						t.Errorf("[%s] logger message is mismatch. %s", tc.l, v["M"])
					}
					if tc.l != v["L"] {
						t.Errorf("[%s] logger level is mismatch. %s", tc.l, v["L"])
					}
				}()
				tc.f(ctx, "%s %d", tc.l, n)
			})
		}
	})
}

func testLoggerAndReader() (*zap.Logger, func(*testing.T) map[string]interface{}) {
	buffer := &bytes.Buffer{}
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(buffer),
		zap.DebugLevel,
	))
	lastLog := func(t *testing.T) map[string]interface{} {
		v := map[string]interface{}{}
		buf, err := io.ReadAll(buffer)
		if err != nil {
			t.Errorf("read log failed")
			return v
		}
		t.Logf("log: %s", string(buf))
		if err := json.Unmarshal(buf, &v); err != nil {
			t.Errorf("decode log failed")
		}
		return v
	}
	return logger, lastLog
}
