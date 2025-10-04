package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime/pprof"

	"github.com/alecthomas/kong"

	"github.com/yeldiRium/3d-rack-brackets/internal/cmd/globals"
	"github.com/yeldiRium/3d-rack-brackets/internal/cmd/render"
)

var cli struct {
	Debug      bool   `help:"Enable debug mode."`
	CPUProfile string `type:"path"`

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

	if err := initializeProfiling(logger); err != nil {
		logger.Error(err.Error())
		ctx.Exit(1)
	}
	defer stopProfiling()

	err := ctx.Run(&globals.Globals{
		Debug:  cli.Debug,
		Logger: logger,
		Stdout: os.Stdout,
	})
	ctx.FatalIfErrorf(err)
}

func initializeProfiling(logger *slog.Logger) error {
	logger.Debug("cpu profiling", slog.Bool("enabled", cli.CPUProfile != ""), slog.String("outPath", cli.CPUProfile))
	if cli.CPUProfile != "" {
		file, err := os.Create(cli.CPUProfile)
		if err != nil {
			return fmt.Errorf("failed to open output file for profiling: %w", err)
		}
		if err := pprof.StartCPUProfile(file); err != nil {
			return fmt.Errorf("failed to initialize profiling: %w", err)
		}
	}

	return nil
}
func stopProfiling() {
	pprof.StopCPUProfile()
}
