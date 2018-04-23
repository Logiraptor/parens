package main

import (
	"testing"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	testCases := map[string]string {
		"list": "()",
		"atoms": "abc def ",
		"atoms 2": "abc def",
		"recursive lists": "(apply (foo bar))",
		"strings": `"string contents"`,
		"ints": `1 2 3`,
		"floats": `1.2 5.0 0.42`,
		"cond": "(cond (= a b) (+ a b) (- a b))",
		"string append": `(append "foo" "bar")`,
		"special forms": "define lambda quote ' quasi-quote ` unquote , cond",
	}
	for name, src := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := lex(newTokenizer([]byte(src)))
			assert.NoError(t, err)
			cupaloy.SnapshotT(t, result)
		})
	}
}
