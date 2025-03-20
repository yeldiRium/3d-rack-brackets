package render

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/yeldiRium/3d-rack-brackets/bin/globals"
)

type RenderCmd struct {
	Output string `arg:"" type:"path" default:"-"`
}

func (render *RenderCmd) Run(globals *globals.Globals) error {
	globals.Logger.Debug("starting to render", slog.Bool("debug", globals.Debug), slog.String("output", render.Output))
	startTime := time.Now()
	defer func() {
		globals.Logger.Debug("done rendering", slog.Duration("elapsed", time.Since(startTime)))
	}()

	var output io.Writer
	if render.Output == "-" {
		output = globals.Stdout
	} else {
		file, err := os.Create(render.Output)
		if err != nil {
			return fmt.Errorf("failed to open output file: %w", err)
		}

		output = file
	}

	fmt.Fprintf(output, "hello :wave:\n")

	return nil
}
