package cc

import (
	"github.com/google/blueprint"

	"github.com/konkers/river"
)

var (
	pctx = blueprint.NewPackageContext("github.com/konkers/river/cc")

	ccCmd = pctx.StaticVariable("ccCmd", "/usr/bin/gcc")

	compile = pctx.StaticRule("compile",
		blueprint.RuleParams{
			Command:     "$ccCmd -o $out $in",
			CommandDeps: []string{"$ccCmd"},
			Description: "Compile $out.",
		})
)

type binary struct {
	properties struct {
		Srcs []string
	}
}

func BinaryFactory() (blueprint.Module, []interface{}) {
	module := new(binary)
	properties := &module.properties
	return module, []interface{}{properties}
}

func (b *binary) GenerateBuildActions(ctx blueprint.ModuleContext) {
	var (
		name       = ctx.ModuleName()
		binaryFile = river.PathForModuleIntermediate(ctx, name)
	)

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:    compile,
		Outputs: []string{binaryFile},
		Inputs:  b.properties.Srcs,
	})
}
