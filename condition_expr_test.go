package ladon

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyRegex(t *testing.T) {
	cases := []string{
		`$hello = {world} and 1=1`,
		`$config startsWith($foo, {bar})`,
		`$earth = {world}`,
		`{ufo} != 'earth'`,
	}

	expectedKeys := []string{
		"world",
		"bar",
		"world",
		"ufo",
	}

	var keys []string
	for _, c := range cases {
		res := keyRegex.FindAllStringSubmatch(c, -1)
		for _, i := range res {
			keys = append(keys, i[1])
		}
	}

	require.Equal(t, expectedKeys, keys)
}

func TestParseExpr(t *testing.T) {
	cases := []string{
		`$hello = $world and 1 = {world}`,
		`$config and startsWith($foo, {bar})`,
		`$config and startsWith($foo, {bar})`,
		`$earth = world`,
	}

	type parseResPair struct {
		hasError      bool
		newExpression string
		values        map[string]interface{}
	}

	expecteds := []parseResPair{
		{
			false,
			`$hello = $world and 1 = 2`,
			map[string]interface{}{
				`world`: 2,
			},
		},
		{
			true,
			``,
			nil,
		},
		{
			false,
			`$config and startsWith($foo, 'bar')`,
			map[string]interface{}{
				`bar`: `bar`,
			},
		},
		{
			false,
			`$earth = world`,
			nil,
		},
	}

	for k, c := range cases {
		expected := expecteds[k]
		expr, err := parseExpr(c, expected.values)
		if expected.hasError {
			require.NotNil(t, err, c)
		} else {
			require.Nil(t, err, c)
		}
		require.Equal(t, expr, expected.newExpression, expr, expected.newExpression)
	}

}

func TestFulfills(t *testing.T) {
	expression := `$filesize <= {pMaxFileSizeLimit} and geoWithin2d($coord2d, {pBottomLeft}, '2,2') and startsWith($field, {pFieldPrefix}) and $_subject = {pSubject}`
	policyVars := map[string]interface{}{
		"pMaxFileSizeLimit": 1000,
		"pBottomLeft":       "0, 0",
		"pTopRight":         "2, 2",
		"pFieldPrefix":      "meera:ac:",
		"pSubject":          "abc",
	}
	ec := ExprCondition{}
	err := ec.Values(expression, policyVars)
	require.Nil(t, err)
	require.Equal(t, `$filesize <= 1000 and geoWithin2d($coord2d, '0, 0', '2,2') and startsWith($field, 'meera:ac:') and $_subject = 'abc'`, ec.Expression)

	ctx := context.Background()
	verifyVars := map[string]interface{}{
		"filesize": 10,
		"coord2d":  "1, 1",
		"field":    "meera:ac:spatial",
	}
	passed := ec.Fulfills(ctx, verifyVars, &Request{Subject: "abc"})
	require.Nil(t, ec.ContextError(), ec.ContextError())
	require.Equal(t, true, passed)

	verifyVars[`filesize`] = 1001
	passed = ec.Fulfills(ctx, verifyVars, &Request{Subject: "abc"})
	require.Nil(t, ec.ContextError(), ec.ContextError())
	require.Equal(t, false, passed)
}
