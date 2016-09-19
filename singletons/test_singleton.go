package singletons

import (
	"github.com/google/blueprint"
	"github.com/konkers/river"
)

type testSingleton struct {
	reports []string
}

var (
	testRunnerCmd = pctx.VariableConfigMethod("testRunnerCmd",
		river.Config.TestRunner)

	runTest = pctx.StaticRule("runTest",
		blueprint.RuleParams{
			Command: "$testRunnerCmd -type $testType -name $testName " +
				"-binary $in > $out",
			CommandDeps: []string{"$testRunnerCmd"},
			Description: "Test $in.",
			Pool:        blueprint.Console,
		}, "testName", "testType")
)

func init() {
	river.RegisterSingletonType("testSingleton", testSingletonFactory)
}

func testSingletonFactory() blueprint.Singleton {
	return new(testSingleton)
}

func (t *testSingleton) generateTestRunner(ctx blueprint.SingletonContext,
	producer river.TestProducer) {
	var (
		config     = ctx.Config().(river.Config)
		testType   = producer.TestType()
		testName   = producer.TestName()
		binaryPath = producer.TestBinaryPath()
		reportPath = config.PathForTestReport(testName)
	)

	t.reports = append(t.reports, reportPath)

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:     runTest,
		Outputs:  []string{reportPath},
		Inputs:   []string{binaryPath},
		Optional: true,
		Args: map[string]string{
			"testName": testName,
			"testType": testType,
		},
	})
}

func (t *testSingleton) GenerateBuildActions(ctx blueprint.SingletonContext) {
	ctx.VisitAllModules(func(module blueprint.Module) {
		if producer, ok := module.(river.TestProducer); ok {
			t.generateTestRunner(ctx, producer)
		}
	})

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:      blueprint.Phony,
		Outputs:   []string{"test"},
		Implicits: t.reports,
		Optional:  true,
	})
}
