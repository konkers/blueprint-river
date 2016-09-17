package cc

import (
	"github.com/google/blueprint"
	"github.com/google/blueprint/pathtools"

	"github.com/konkers/river"
)

var (
	pctx = blueprint.NewPackageContext("github.com/konkers/river/cc")

	ccCmd = pctx.StaticVariable("ccCmd", "/usr/bin/gcc")

	compile = pctx.StaticRule("compile",
		blueprint.RuleParams{
			Command:     "$ccCmd -c -o $out $in",
			CommandDeps: []string{"$ccCmd"},
			Description: "Compile $out.",
		})

	link = pctx.StaticRule("link",
		blueprint.RuleParams{
			Command:     "$ccCmd -o $out $in",
			CommandDeps: []string{"$ccCmd"},
			Description: "Link $out.",
		})
)

func init() {
	river.RegisterModuleType("cc_binary", binaryFactory)
}

type binary struct {
	properties struct {
		Srcs []string
	}
}

func binaryFactory() (blueprint.Module, []interface{}) {
	module := new(binary)
	properties := &module.properties
	return module, []interface{}{properties}
}

func (b *binary) GenerateBuildActions(ctx blueprint.ModuleContext) {
	var (
		name       = ctx.ModuleName()
		binaryFile = river.PathForModuleIntermediate(ctx, name)
	)

	objFiles := make([]string, len(b.properties.Srcs))

	for _, src := range b.properties.Srcs {
		objName := pathtools.ReplaceExtension(src, "o")
		objFile := river.PathForModuleIntermediate(ctx, objName)
		objFiles = append(objFiles, objFile)

		ctx.Build(pctx, blueprint.BuildParams{
			Rule:    compile,
			Outputs: []string{objFile},
			Inputs:  b.properties.Srcs,
		})
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:    link,
		Outputs: []string{binaryFile},
		Inputs:  objFiles,
	})
}
