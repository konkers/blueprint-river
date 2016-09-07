package river

import (
	"path/filepath"

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
		filepath.Join(paths...))
}
