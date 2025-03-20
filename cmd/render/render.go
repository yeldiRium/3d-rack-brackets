package render

import (
	"fmt"

	"github.com/yeldiRium/3d-rack-brackets/bin/globals"
)

type RenderCmd struct {
	Output string `arg:"" type:"path"`
}

func (render *RenderCmd) Run(globals *globals.Globals) error {
	fmt.Printf("debugging: %v\n", globals.Debug)
	fmt.Printf("rendering to %v\n", render.Output)

	return nil
}
