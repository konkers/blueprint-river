// Based on bootstrap go rules from blueprint:
// Copyright 2014 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package golang

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/blueprint"

	"github.com/konkers/river"
)

var (
	pctx = blueprint.NewPackageContext("github.com/konkers/river/golang")

	hostPrebuiltTag = pctx.VariableConfigMethod("hostPrebuiltTag",
		river.Config.HostPrebuiltTag)
	goRoot = pctx.StaticVariable("goRoot",
		filepath.Join(".", "prebuilts", "go", "$hostPrebuiltTag"))
	compileCmd = pctx.StaticVariable("compileCmd", toolPath("compile"))
	linkCmd    = pctx.StaticVariable("linkCmd", toolPath("link"))

	compile = pctx.StaticRule("compile",
		blueprint.RuleParams{
			Command: "GOROOT='$goRoot' $compileCmd -o $out -p $pkgPath " +
				"-complete $incFlags -pack $in",
			CommandDeps: []string{"$compileCmd"},
			Description: "compile $out",
		},
		"incFlags", "pkgPath")

	link = pctx.StaticRule("link",
		blueprint.RuleParams{
			Command:     "GOROOT='$goRoot' $linkCmd -o $out $libDirFlags $in",
			CommandDeps: []string{"$linkCmd"},
			Description: "link $out",
		},
		"libDirFlags")
)

type common struct {
	properties struct {
		Srcs    []string
		PkgPath string
	}

	pkgRoot     string
	archiveFile string
}

type binary struct {
	common

	binaryFile string
}

type pkg struct {
	common
}

func init() {
	river.RegisterModuleType("go_binary", binaryFactory)
	river.RegisterModuleType("go_package", packageFactory)
}

func binaryFactory() (blueprint.Module, []interface{}) {
	b := new(binary)
	return b, []interface{}{&b.common.properties}
}

func packageFactory() (blueprint.Module, []interface{}) {
	p := new(pkg)
	return p, []interface{}{&p.common.properties}
}

func (c *common) GenerateBuildActions(ctx blueprint.ModuleContext) {
	c.pkgRoot = river.PathForModuleIntermediate(ctx, "pkg")
	c.archiveFile = filepath.Join(c.pkgRoot,
		filepath.FromSlash(c.properties.PkgPath)+".a")

	srcFiles := make([]string, len(c.properties.Srcs))
	for _, src := range c.properties.Srcs {
		srcFiles = append(srcFiles, river.PathForModuleSource(ctx, src))
	}

	var incFlags []string
	var deps []string
	ctx.VisitDepsDepthFirstIf(isGoPackageProducer,
		func(module blueprint.Module) {
			dep := module.(goPackageProducer)
			incDir := dep.GoPkgRoot()
			target := dep.GoPackageTarget()
			incFlags = append(incFlags, "-I "+incDir)
			deps = append(deps, target)
		})

	compileArgs := map[string]string{
		"pkgPath": c.properties.PkgPath,
	}

	if len(incFlags) > 0 {
		compileArgs["incFlags"] = strings.Join(incFlags, " ")
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:      compile,
		Outputs:   []string{c.archiveFile},
		Inputs:    srcFiles,
		Implicits: deps,
		Args:      compileArgs,
	})
}

func (b *binary) GenerateBuildActions(ctx blueprint.ModuleContext) {
	var (
		name = ctx.ModuleName()
	)

	b.binaryFile = river.PathForModuleIntermediate(ctx, name)
	b.common.GenerateBuildActions(ctx)

	var libDirFlags []string
	ctx.VisitDepsDepthFirstIf(isGoPackageProducer,
		func(module blueprint.Module) {
			dep := module.(goPackageProducer)
			libDir := dep.GoPkgRoot()
			libDirFlags = append(libDirFlags, "-L "+libDir)
		})

	linkArgs := map[string]string{}
	if len(libDirFlags) > 0 {
		linkArgs["libDirFlags"] = strings.Join(libDirFlags, " ")
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:    link,
		Outputs: []string{b.binaryFile},
		Inputs:  []string{b.archiveFile},
		Args:    linkArgs,
	})
}

func (p *pkg) GenerateBuildActions(ctx blueprint.ModuleContext) {
	p.common.GenerateBuildActions(ctx)
}

func (p *pkg) LibraryFileName() string {
	return p.archiveFile
}

func toolPath(tool string) string {
	goTag := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
	return filepath.Join("$goRoot", "pkg", "tool", goTag, tool)
}

type goPackageProducer interface {
	GoPkgRoot() string
	GoPackageTarget() string
}

func isGoPackageProducer(module blueprint.Module) bool {
	_, ok := module.(goPackageProducer)
	return ok
}

func (p *pkg) GoPkgRoot() string {
	return p.pkgRoot
}

func (p *pkg) GoPackageTarget() string {
	return p.archiveFile
}
