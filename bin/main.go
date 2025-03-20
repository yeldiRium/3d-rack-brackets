package main

import (
	"github.com/alecthomas/kong"

	"github.com/yeldiRium/3d-rack-brackets/bin/globals"
	"github.com/yeldiRium/3d-rack-brackets/cmd/render"
)

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Render render.RenderCmd `cmd:"" help:"render the rack"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.UsageOnError(),
	)
	err := ctx.Run(&globals.Globals{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
