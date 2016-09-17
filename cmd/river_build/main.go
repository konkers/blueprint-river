package main

import (
	"flag"
	"path/filepath"

	"github.com/google/blueprint/bootstrap"

	"github.com/konkers/river"
)

func main() {
	flag.Parse()

	// The top-level Blueprints file is passed as the first argument.
	srcDir := filepath.Dir(flag.Arg(0))

	// Create the build context.
	ctx := river.NewContext()

	// Create and initialize the custom Config object.
	config, err := river.NewConfig(srcDir, bootstrap.BuildDir)
	if err != nil {
		panic(err)
	}

	// This call never returns
	bootstrap.Main(ctx, config)
}
