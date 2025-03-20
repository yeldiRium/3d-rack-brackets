package render

import (
	"fmt"
	"io"
	"os"

	"github.com/yeldiRium/3d-rack-brackets/bin/globals"
)

type RenderCmd struct {
	Output string `arg:"" type:"path" default:"-"`
}

func (render *RenderCmd) Run(globals *globals.Globals) error {
	fmt.Printf("debugging: %v\n", globals.Debug)
	fmt.Printf("rendering to %v\n", render.Output)

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

	fmt.Fprintf(output, "hello :wave:")

	return nil
}
