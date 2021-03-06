package cc

import (
	"path/filepath"
	"strings"

	"github.com/google/blueprint"
	"github.com/google/blueprint/pathtools"

	"github.com/konkers/river"
)

var (
	pctx = blueprint.NewPackageContext("github.com/konkers/river/cc")

	// TODO(konkers): Make this configurable.
	testMainDep = "c_test"

	// TODO(konkers): Replace with host/target specific config.
	hostPrebuiltTag = pctx.VariableConfigMethod("hostPrebuiltTag",
		river.Config.HostPrebuiltTag)
	ccCmd = pctx.StaticVariable("ccCmd",
		"./prebuilts/clang/$hostPrebuiltTag/bin/clang")
	arCmd = pctx.StaticVariable("arCmd", "/usr/bin/ar")

	// TODO(konkers): Implement include sandboxing.
	cFlags = pctx.StaticVariable("cFlags", "-I. -O3 -Wall -Werror")

	compile = pctx.StaticRule("compile",
		blueprint.RuleParams{
			Command: "$ccCmd -cc1 $cFlags $extraCFlags " +
				"-MT $out -dependency-file ${out}.d " +
				"-emit-obj -o $out $in",
			CommandDeps: []string{"$ccCmd"},
			Description: "Compile $out.",
			Depfile:     "${out}.d",
			Deps:        blueprint.DepsGCC,
		}, "extraCFlags")

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
	river.RegisterModuleType("cc_test", testFactory)

	river.RegisterBottomUpMutator("testDepsMutator", testDepsMutator)
}

type common struct {
	properties struct {
		Srcs        []string
		ExtraCFlags []string
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

type test struct {
	binary

	name string
}

func binaryFactory() (blueprint.Module, []interface{}) {
	b := new(binary)
	return b, []interface{}{&b.common.properties}
}

func libraryFactory() (blueprint.Module, []interface{}) {
	l := new(library)
	return l, []interface{}{&l.common.properties, &l.properties}
}

func testFactory() (blueprint.Module, []interface{}) {
	t := new(test)
	return t, []interface{}{&t.common.properties}
}

func (c *common) GenerateBuildActions(ctx blueprint.ModuleContext) {
	var (
		extraCFlags = strings.Join(c.properties.ExtraCFlags, " ")
	)

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
			Args: map[string]string{
				"extraCFlags": extraCFlags,
			},
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

func (t *test) GenerateBuildActions(ctx blueprint.ModuleContext) {
	t.binary.GenerateBuildActions(ctx)
	t.name = filepath.Join(ctx.ModuleDir(), ctx.ModuleName())
}

func (t *test) TestType() string {
	return "c_test"
}

func (t *test) TestName() string {
	return t.name
}

func (t *test) TestBinaryPath() string {
	return t.binaryFile
}

func testDepsMutator(ctx blueprint.BottomUpMutatorContext) {
	if _, ok := ctx.Module().(*test); ok {
		ctx.AddDependency(ctx.Module(), nil, testMainDep)
	}
}
