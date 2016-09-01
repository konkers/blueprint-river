package main

import (
	"flag"
	_ "path/filepath"

	"github.com/google/blueprint"
	"github.com/google/blueprint/bootstrap"
)

func main() {
	flag.Parse()

	// The top-level Blueprints file is passed as the first argument.
	//	srcDir := filepath.Dir(flag.Arg(0))

	// Create the build context.
	ctx := blueprint.NewContext()

	// Register custom module types
	//ctx.RegisterModuleType("foo", logic.FooModule)
	//ctx.RegisterModuleType("bar", logic.BarModule)

	// Register custom singletons
	//ctx.RegisterSingleton("baz", logic.NewBazSingleton())

	// Create and initialize the custom Config object.
	//config := logic.NewConfig(srcDir)

	// This call never returns
	bootstrap.Main(ctx, nil)
}
