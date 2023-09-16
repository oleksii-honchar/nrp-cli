package blablo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"

	"beaver/blablo/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	logger *log.Logger
}

func (handler *PrettyHandler) Handle(ctx context.Context, rec slog.Record) error {
	level := rec.Level.String()

	switch rec.Level {
	case slog.LevelDebug:
		level = color.Magenta + level + color.Reset
	case slog.LevelInfo:
		level = color.Blue + level + color.Reset
	case slog.LevelWarn:
		level = color.Yellow + level + color.Reset
	case slog.LevelError:
		level = color.Red + level + color.Reset
	}
	level += " |"

	fields := make(map[string]interface{}, rec.NumAttrs())
	rec.Attrs(func(attr slog.Attr) bool {
		fields[attr.Key] = attr.Value.Any()

		return true
	})

	jsonStr, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	// timeStr := rec.Time.Format("2006/01/02 15:05:05.000")
	msg := color.White + rec.Message + color.Reset

	// line := fmt.Sprintf("%s %s %s", color.Gray247+timeStr+color.Reset, level, msg)
	line := fmt.Sprintf("%s %s", level, msg)

	if len(fields) > 0 {
		line += " " + color.Yellow + string(jsonStr) + color.Reset
	}

	handler.logger.Println(line)

	return nil
}

func NewPrettyHandler(
	prefix string,
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	normPrefix := TrimOrPadStringRight(prefix, 10) + " | "

	handler := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		logger:  log.New(out, normPrefix, 0),
	}

	return handler
}
