package ladon

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/EchoUtopia/expr"
)

const (
	ExprConditionName = `ExprCondition`
)

var (
	keyRegex = regexp.MustCompile(`{([_a-zA-Z][_a-zA-Z0-9]*)}`)
)

type ExprValues map[string]interface{}

type ExprCondition struct {
	// expression can contains variables, it's format is $***
	// eg: $filesize lte $pMaxFileSizeLimit and geo.within2d($coord2d, $pBottomLeft, $pTopRight)
	//   and startswith($field, $pFieldPrefix) and $subject eq $pSubject
	Expression string `json:"expression"`

	// unlike the other init-at-beginning & read-only-later fields,
	// the following fields are internal rw state, which means thread unsafe.
	err error
}

func (ec *ExprCondition) Fulfills(ctx context.Context, values interface{}, r *Request) bool {

	vars, ok := values.(map[string]interface{})
	if !ok {
		ec.err = errors.New("invalid context value, should setted the value for the varibale of expression condition")
		return false
	}
	if _, ok := vars["_subject"]; !ok {
		vars["_subject"] = r.Subject
	}
	if _, ok := vars["_resource"]; !ok {
		vars["_resource"] = r.Resource
	}
	if _, ok := vars["_action"]; !ok {
		vars["_action"] = r.Action
	}

	result, err := expr.Evaluate(ec.Expression, vars, nil)
	if err != nil {
		ec.err = err
		return false
	}
	return result
}

// Value set value to the condition's Value.
func (ec *ExprCondition) Values(expression string, values map[string]interface{}) error {
	newExpr, err := parseExpr(expression, values)
	if err != nil {
		return err
	}
	ec.Expression = newExpr
	return nil
}

// ContextError returns error in the request context.
func (ec *ExprCondition) ContextError() error {
	return ec.err
}

// GetName returns the condition's name.
func (ec *ExprCondition) GetName() string {
	return ExprConditionName
}

func quoteForExpr(str string) string {
	str = strings.ReplaceAll(str, `\`, `\\`)
	str = strings.ReplaceAll(str, `'`, `\'`)
	return `'` + str + `'`
}

// parseExpr replace placeholders like: '{***}'  with template recognised token, eg {{ .var }} in the input expr,
// and check if all placeholders feed with values, and then put the parsed tree into parseTreePool
func parseExpr(expression string, values map[string]interface{}) (string, error) {

	expectedKeys := []string{}
	expression = keyRegex.ReplaceAllStringFunc(expression, func(s string) string {
		res := keyRegex.FindAllStringSubmatch(s, -1)
		if len(res) == 0 {
			return s
		}
		matchedKey := res[0][1]
		expectedKeys = append(expectedKeys, matchedKey)
		return fmt.Sprintf("{{ .%s }}", matchedKey)
	})
	for _, k := range expectedKeys {
		iv, ok := values[k]
		if !ok {
			return ``, fmt.Errorf(`value of placeholder: "%s" not provided`, k)
		}
		v, ok := iv.(string)
		if ok {
			values[k] = quoteForExpr(v)
		}
	}

	buf := &bytes.Buffer{}
	tpl := template.New("")
	tpl, err := tpl.Parse(expression)
	if err != nil {
		return ``, err
	}
	if err := tpl.Execute(buf, values); err != nil {
		return ``, err
	}
	expression = buf.String()
	parser := expr.NewParser()
	if _, err := parser.ParseWithCache(expression); err != nil {
		return ``, fmt.Errorf("%s, %w", expression, err)
	}
	return expression, nil
}
