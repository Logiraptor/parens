package main

import (
	"testing"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	testCases := map[string]string {
		"nil": "()",
		"atom": "abc",
		"condition": "(cond ((= a b) (+ a b)) ((neq a b) (- a b)) )",
		"arithmetic": "(+ 1 2)",
		"string append": `(append "foo" "bar")`,
		"float addition": `(+ 1.5 0.34)`,
		"lambda": `(lambda (x) (+ x 1))`,
		"definition": `(define a 1)`,
		"quoted forms": `(define q '(1 2) (quote 1))`,
		"quasi-quoted forms": "(define q `(,1 ,2) `(unquote (quasi-quote 1)))",
	}
	for name, src := range testCases {
		t.Run(name, func(t *testing.T) {
			result, err := lex(newTokenizer([]byte(src)))
			assert.NoError(t, err)
			ast, err := parse(newAstEmitter(result))
			assert.NoError(t, err)
			cupaloy.SnapshotT(t, ast)
		})
	}
}
