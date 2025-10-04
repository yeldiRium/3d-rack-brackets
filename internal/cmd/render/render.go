package render

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/primitive"

	"github.com/yeldiRium/3d-rack-brackets/internal/cmd/globals"
	"github.com/yeldiRium/3d-rack-brackets/internal/ghostscad"
	"github.com/yeldiRium/3d-rack-brackets/internal/shapes"
	"github.com/yeldiRium/3d-rack-brackets/internal/shapes/rack"
)

type RenderCmd struct {
	Production bool   `short:"p"`
	Output     string `arg:""    default:"-" type:"path"`
}

func (render *RenderCmd) Run(globals *globals.Globals) error {
	globals.Logger.Debug("starting to render", slog.Bool("debug", globals.Debug), slog.String("output", render.Output))
	startTime := time.Now()
	defer func() {
		globals.Logger.Debug("done rendering", slog.Duration("elapsed", time.Since(startTime)))
	}()

	if render.Production {
		ghostscad.SetFa(5)
		ghostscad.SetFs(0.5)
	}

	output, err := render.ChooseOutput(globals.Stdout)
	if err != nil {
		return fmt.Errorf("failed to open output stream: %w", err)
	}
	bufferedOutput := bufio.NewWriter(output)

	segmentCount := 3
	shape := rack.MakeRack(uint8(segmentCount))
	err = shapes.ResolveAnchors(shape.Foot)
	if err != nil {
		return fmt.Errorf("failed to resolve anchors: %w", err)
	}

	orientedShape := primitive.NewRotation(mgl64.Vec3{0, 0, 0}, shape)
	translatedShape := primitive.NewTranslation(mgl64.Vec3{0, 0, float64(rack.RackFootThicknessFront)}, orientedShape)

	ghostscad.RenderGlobals(bufferedOutput)
	translatedShape.Render(bufferedOutput)

	err = bufferedOutput.Flush()
	if err != nil {
		return err
	}

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
