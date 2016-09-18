package river

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/google/blueprint"
)

type Config struct {
	*config
}

type config struct {
	srcDir   string
	buildDir string
}

func NewConfig(srcDir, buildDir string) (Config, error) {
	config := &config{
		srcDir:   srcDir,
		buildDir: buildDir,
	}

	return Config{config}, nil
}

func PathForModuleIntermediate(ctx blueprint.ModuleContext, paths ...string) string {
	config := ctx.Config().(Config)
	return filepath.Join(config.buildDir, "intermediates", ctx.ModuleDir(),
		ctx.ModuleName(), filepath.Join(paths...))
}

func PathForModuleSource(ctx blueprint.ModuleContext, paths ...string) string {
	config := ctx.Config().(Config)
	return filepath.Join(config.srcDir, ctx.ModuleDir(),
		filepath.Join(paths...))
}

// Returns the prebuilt tag for the host (i.e darwin-x86_64).
func (c *config) HostPrebuiltTag() string {
	var tag string

	switch runtime.GOOS {
	case "linux":
		tag = "linux-"
	case "darwin":
		tag = "darwin-"
	default:
		panic(fmt.Sprintf("Unknown host OS %s", runtime.GOOS))
	}

	switch runtime.GOARCH {
	case "amd64":
		tag += "x86_64"
	default:
		panic(fmt.Sprintf("Unknown host arch %s", runtime.GOARCH))
	}

	return tag
}

// Returns the root path of the specified prebuilt tool.
func (c *config) HostPrebuiltRoot(tool string) string {
	return filepath.Join(c.srcDir, "prebuilts", c.HostPrebuiltTag())
}
