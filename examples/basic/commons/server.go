package example

import (
	"log"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spf13/afero"
	"github.com/terraform-linters/tflint/tflint"
	"github.com/zclconf/go-cty/cty"
)

type Server struct {
	Runner *tflint.Runner
}

func (s *Server) Walk(args interface{}, resp *hcl.Attribute) error {
	loader, err := tflint.NewLoader(afero.Afero{Fs: afero.NewOsFs()}, tflint.EmptyConfig())
	if err != nil {
		panic(err)
	}
	cfg, err := loader.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	s.Runner, err = tflint.NewRunner(tflint.EmptyConfig(), map[string]tflint.Annotations{}, cfg)
	if err != nil {
		panic(err)
	}

	err = s.Runner.WalkResourceAttributes("aws_instance", "instance_type", func(attr *hcl.Attribute) error {
		log.Printf("hcl.Attribute: %s", attr)
		*resp = *attr
		return nil
	})
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *Server) EvaluateExpr(expr hcl.Expression, resp *cty.Value) error {
	ctx := terraform.BuiltinEvalContext{
		Evaluator: &terraform.Evaluator{
			Config:             s.Runner.TFConfig,
			VariableValues:     prepareVariableValues(s.Runner.TFConfig.Module.Variables),
			VariableValuesLock: &sync.Mutex{},
		},
	}

	val, diags := ctx.EvaluateExpr(expr, cty.DynamicPseudoType, nil)
	if diags.HasErrors() {
		panic(diags)
	}
	log.Printf("cty.Value: %s", val)
	*resp = val

	return nil
}

func prepareVariableValues(configVars map[string]*configs.Variable) map[string]map[string]cty.Value {
	overrideVariables := terraform.DefaultVariableValues(configVars)

	variableValues := make(map[string]map[string]cty.Value)
	variableValues[""] = make(map[string]cty.Value)
	for k, iv := range overrideVariables {
		variableValues[""][k] = iv.Value
	}
	return variableValues
}
