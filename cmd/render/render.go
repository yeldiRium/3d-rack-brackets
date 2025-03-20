package render

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/yeldiRium/3d-rack-brackets/bin/globals"
	"github.com/yeldiRium/3d-rack-brackets/ghostscad"
	"github.com/yeldiRium/3d-rack-brackets/shapes/rack"
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

	output, err := render.ChooseOutput(globals.Stdout)
	if err != nil {
		return fmt.Errorf("failed to open output stream: %w", err)
	}
	bufferedOutput := bufio.NewWriter(output)

	shape := rack.MakeRack()

	ghostscad.RenderGlobals(bufferedOutput)
	shape.Render(bufferedOutput)

	bufferedOutput.Flush()

	return nil
}

func (render *RenderCmd) ChooseOutput(stdout io.Writer) (io.Writer, error) {
	var output io.Writer
	if render.Output == "-" {
		output = stdout
	} else {
		file, err := os.Create(render.Output)
		if err != nil {
			return nil, fmt.Errorf("failed to open output file: %w", err)
		}

		output = file
	}

	return output, nil
}
