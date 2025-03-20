package main

import (
	"log/slog"
	"os"

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

	logLevel := slog.LevelInfo
	if cli.Debug {
		logLevel = slog.LevelDebug
	}		
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)

	err := ctx.Run(&globals.Globals{
		Debug:  cli.Debug,
		Logger: logger,
		Stdout: os.Stdout,
	})
	ctx.FatalIfErrorf(err)
}
