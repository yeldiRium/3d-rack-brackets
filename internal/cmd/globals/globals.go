package globals

import (
	"io"
	"log/slog"
)

type Globals struct {
	Debug  bool
	Logger *slog.Logger
	Stdout io.Writer
}
