/*
 * @Author: yujiajie
 * @Date: 2025-07-11 17:05:38
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-07-11 17:05:53
 * @FilePath: /Go-Base/pkg/logger/core.go
 * @Description:
 */
package logger

import (
	"bytes"
	"encoding/json"
	"fmt"

	"go.uber.org/zap/zapcore"
)

type JsonCore struct {
	zapcore.Core
	enc zapcore.Encoder
	ws  zapcore.WriteSyncer
	lvl zapcore.LevelEnabler
}

func NewJsonCore(enc zapcore.Encoder, ws zapcore.WriteSyncer, lvl zapcore.LevelEnabler) *JsonCore {
	return &JsonCore{
		enc: enc,
		ws:  ws,
		lvl: lvl,
	}
}

func (c *JsonCore) With(fields []zapcore.Field) zapcore.Core {
	return &JsonCore{
		enc: c.enc.Clone(),
		ws:  c.ws,
		lvl: c.lvl,
	}
}

func (c *JsonCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *JsonCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	var msgJson map[string]any
	if err := json.Unmarshal([]byte(ent.Message), &msgJson); err != nil {
		//如果不是合法json，直接写字符串
		c.ws.Write([]byte(fmt.Sprintf(`{"raw_msg":%q}\n`, ent.Message)))
		return err
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(msgJson); err != nil {
		return err
	}

	_, err := c.ws.Write(buf.Bytes())
	return err
}

func (c *JsonCore) Enabled(level zapcore.Level) bool {
	return c.lvl.Enabled(level)
}

func (c *JsonCore) Sync() error {
	return c.ws.Sync()
}
