package cc

import (
	"github.com/google/blueprint"
	"github.com/google/blueprint/pathtools"

	"github.com/konkers/river"
)

var (
	pctx = blueprint.NewPackageContext("github.com/konkers/river/cc")

	// TODO(konkers): Replace with host/target specific config.
	ccCmd = pctx.StaticVariable("ccCmd", "/usr/bin/gcc")
	arCmd = pctx.StaticVariable("arCmd", "/usr/bin/ar")

	// TODO(konkers): Implement include sandboxing.
	cFlags = pctx.StaticVariable("cFlags", "-I.")

	compile = pctx.StaticRule("compile",
		blueprint.RuleParams{
			Command:     "$ccCmd $cFlags -c -o $out $in",
			CommandDeps: []string{"$ccCmd"},
			Description: "Compile $out.",
		})

	link = pctx.StaticRule("link",
		blueprint.RuleParams{
			Command:     "$ccCmd -o $out $in",
			CommandDeps: []string{"$ccCmd"},
			Description: "Link $out.",
		})

	staticLib = pctx.StaticRule("ar",
		blueprint.RuleParams{
			Command:     "$arCmd cr $out $in",
			CommandDeps: []string{"$arCmd"},
			Description: "Static Lib $out.",
		})
)

func init() {
	river.RegisterModuleType("cc_binary", binaryFactory)
	river.RegisterModuleType("cc_library", libraryFactory)
}

type common struct {
	properties struct {
		Srcs []string
	}

	objFiles []string
}

type binary struct {
	common

	binaryFile string
	linkFlags  string
}

type library struct {
	common

	properties struct {
		Incs []string
	}

	libraryFile string
}

func binaryFactory() (blueprint.Module, []interface{}) {
	b := new(binary)
	return b, []interface{}{&b.common.properties}
}

func libraryFactory() (blueprint.Module, []interface{}) {
	l := new(library)
	return l, []interface{}{&l.common.properties, &l.properties}
}

func (c *common) GenerateBuildActions(ctx blueprint.ModuleContext) {
	c.objFiles = make([]string, len(c.properties.Srcs))

	for _, src := range c.properties.Srcs {
		objName := pathtools.ReplaceExtension(src, "o")
		objFile := river.PathForModuleIntermediate(ctx, objName)
		srcFile := river.PathForModuleSource(ctx, src)
		c.objFiles = append(c.objFiles, objFile)

		ctx.Build(pctx, blueprint.BuildParams{
			Rule:    compile,
			Outputs: []string{objFile},
			Inputs:  []string{srcFile},
		})
	}
}

func (b *binary) GenerateBuildActions(ctx blueprint.ModuleContext) {
	var (
		name = ctx.ModuleName()
	)

	b.binaryFile = river.PathForModuleIntermediate(ctx, name)
	b.common.GenerateBuildActions(ctx)

	inputs := b.objFiles
	ctx.VisitDepsDepthFirstIf(river.IsLibraryProducer,
		func(module blueprint.Module) {
			libProducer := module.(river.LibraryProducer)
			inputs = append(inputs, libProducer.LibraryFileName())
		})

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:    link,
		Outputs: []string{b.binaryFile},
		Inputs:  inputs,
	})
}

func (l *library) GenerateBuildActions(ctx blueprint.ModuleContext) {
	var (
		name = ctx.ModuleName()
	)

	l.libraryFile = river.PathForModuleIntermediate(ctx, name) + ".a"

	l.common.GenerateBuildActions(ctx)

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:    staticLib,
		Outputs: []string{l.libraryFile},
		Inputs:  l.objFiles,
	})
}

func (l *library) LibraryFileName() string {
	return l.libraryFile
}
