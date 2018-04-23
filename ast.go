package main

import "strconv"

//go:generate stringer -type=astType

type astType int

type ast interface {
	isAst()
}

type atom struct {
	value string
}

func (atom) isAst() {}

func newAtom(s string) ast {
	return atom{s}
}

type intLiteral struct {
	value int
}

func (intLiteral) isAst() {}

func newInteger(s string) (ast, error) {
	val, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}
	return intLiteral{value: val}, nil
}

type floatLiteral struct {
	value float64
}

func (floatLiteral) isAst() {}

func newFloat(s string) (ast, error) {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}
	return floatLiteral{value: val}, nil
}

type stringLiteral struct {
	value string
}

func (stringLiteral) isAst() {}

func newString(s string) (ast) {
	return stringLiteral{value: s[1:len(s)-1],}
}

type list struct {
	list []ast
}

func (list) isAst() {}

func newList(elements []ast) ast {
	if len(elements) > 0 {
		if atomNode, ok := elements[0].(atom); ok {
			switch atomNode.value {
			case "define":
				name := elements[1].(atom).value
				return newDefine(name, elements[2])
			case "lambda":
				args := elements[1]
				body := elements[2]
				lmda, _ := newLambda(args, body)
				return  lmda
			case "quote":
				form := elements[1]
				return newQuote(form)
			case "quasi-quote":
				form := elements[1]
				return newQuasiQuote(form)
			case "unquote":
				form := elements[1]
				return newUnquote(form)
			case "cond":
				conditionForm := newCond()
				for _, elem := range elements[1:] {
					condition := elem.(list).list[0]
					conseq := elem.(list).list[1]
					conditionForm.clauses = append(conditionForm.clauses, condClause{
						condition: condition,
						consequence: conseq,
					})
				}
				return conditionForm
			}
		}
	}

	return list{list: elements,}
}

type define struct {
	name  string
	value ast
}

func (define) isAst() {}

func newDefine(name string, value ast) ast {
	return define{name, value}
}

type lambda struct {
	args ast
	body ast
}

func (lambda) isAst() {}

func newLambda(args, body ast) (ast, error) {
	return lambda{args, body}, nil
}

type condClause struct {
	condition   ast
	consequence ast
}

type cond struct {
	clauses []condClause
}

func (cond) isAst() {}

func newCond() (cond) {
	return cond{}
}

type quote struct {
	form ast
}

func (quote) isAst() {}

func newQuote(form ast) (ast) {
	return quote{form}
}

type unquote struct {
	form ast
}

func (unquote) isAst() {}

func newUnquote(form ast) (ast) {
	return unquote{form}
}

type quasiQuote struct {
	form ast
}

func (quasiQuote) isAst() {}

func newQuasiQuote(form ast) (ast) {
	return quasiQuote{form}
}
