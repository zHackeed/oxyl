package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

// Todo: reshape to be able to send context data before messages

var (
	InfoColor  = color.New(color.Bold, color.FgGreen)
	WarnColor  = color.New(color.Bold, color.FgYellow)
	ErrorColor = color.New(color.Bold, color.FgRed)
	DebugColor = color.New(color.Bold, color.FgBlue)

	DateColor = color.New(color.FgHiBlack, color.Italic)

	FieldKeyColor   = color.New(color.FgHiCyan)
	FieldValueColor = color.New(color.FgHiWhite)
	FieldSeparator  = color.New(color.FgHiBlack)

	MessageColor = color.New(color.FgHiWhite, color.Bold)
)

var levelLabels = map[slog.Level]string{
	slog.LevelDebug: DebugColor.Sprint("DBG"),
	slog.LevelInfo:  InfoColor.Sprint("INF"),
	slog.LevelWarn:  WarnColor.Sprint("WRN"),
	slog.LevelError: ErrorColor.Sprint("ERR"),
}

// Ref impl logic: https://betterstack.com/community/guides/logging/logging-in-go/

type PrettyHandler struct {
	slog.Handler
	w io.Writer
}

func NewPrettyHandler(w io.Writer, opts slog.HandlerOptions) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewTextHandler(w, &opts),
		w:       w,
	}
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level, ok := levelLabels[r.Level]
	if !ok {
		level = r.Level.String()
	}

	sep := FieldSeparator.Sprint("│")

	var b strings.Builder
	b.WriteString(FieldSeparator.Sprint("["))
	b.WriteString(DateColor.Sprint(r.Time.Format(time.RFC3339)))
	b.WriteString(FieldSeparator.Sprint("]"))
	b.WriteByte(' ')
	b.WriteString(level)
	b.WriteByte(' ')
	b.WriteString(sep)
	b.WriteByte(' ')
	b.WriteString(MessageColor.Sprint(r.Message))

	r.Attrs(func(attr slog.Attr) bool {
		b.WriteByte(' ')
		b.WriteString(sep)
		b.WriteByte(' ')

		// Shallow the bad key 'key=value'. Looks ugly honestly.
		if attr.Key == "!BADKEY" {
			b.WriteString(FieldValueColor.Sprint(attr.Value.String()))
		} else {
			b.WriteString(FieldKeyColor.Sprint(attr.Key))
			b.WriteString(FieldSeparator.Sprint(" → "))
			b.WriteString(FieldValueColor.Sprint(attr.Value.String()))
		}

		return true
	})

	b.WriteByte('\n') // new line
	_, err := fmt.Fprint(h.w, b.String())
	return err
}

// Register ->> Avoid duplicate code
func Register(ops slog.HandlerOptions) {
	customLogger := NewPrettyHandler(os.Stdout, ops)
	slog.SetDefault(slog.New(customLogger))
}
