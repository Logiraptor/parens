package main

import (
	"io/ioutil"
	"fmt"
	"os"
)

type module struct {
}

func main() {
	buf, err := ioutil.ReadFile("./simple.parens")
	if err != nil {
		fail(err)
	}
	fmt.Println(buf)
	tokens, err := lex(newTokenizer(buf))
	if err != nil {
		fail(err)
	}
	fmt.Println(tokens)
	ast, err := parse(newAstEmitter(tokens))
	if err != nil {
		fail(err)
	}
	fmt.Println(ast)
	module, err := load(ast)
	if err != nil {
		fail(err)
	}
	fmt.Println(module)
	result, err := execute(module)
	if err != nil {
		fail(err)
	}
	fmt.Println(result)

	fmt.Println(callFuncDynamically("strings", "Index", "abc", "a"))

}

func fail(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func execute(m module) (int, error) {
	return 0, nil
}

func load(a ast) (module, error) {
	return module{}, nil
}
