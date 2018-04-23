package main

import (
	"fmt"
)

type astEmitter struct {
	tokens []token
	pos    int
}

func newAstEmitter(tokens []token) *astEmitter {
	return &astEmitter{
		tokens: tokens,
	}
}

func (a *astEmitter) shift() {
	a.pos++
}

func (a *astEmitter) curr() token {
	return a.tokens[a.pos]
}

func parse(emitter *astEmitter) (ast, error) {
	switch emitter.curr().typ {
	case QuoteTok:
		emitter.shift()
		quotedForm, err := parse(emitter)
		if err != nil {
			return nil, err
		}
		return newQuote(quotedForm), nil
	case UnquoteTok:
		emitter.shift()
		quotedForm, err := parse(emitter)
		if err != nil {
			return nil, err
		}
		return newUnquote(quotedForm), nil
	case QuasiQuoteTok:
		emitter.shift()
		quotedForm, err := parse(emitter)
		if err != nil {
			return nil, err
		}
		return newQuasiQuote(quotedForm), nil
	case AtomTok:
		atom := newAtom(string(emitter.curr().buf))
		emitter.shift()
		return atom, nil
	case IntTok:
		integer, err := newInteger(string(emitter.curr().buf))
		if err != nil {
			return nil, err
		}
		emitter.shift()
		return integer, nil
	case FloatTok:
		float, err := newFloat(string(emitter.curr().buf))
		if err != nil {
			return nil, err
		}
		emitter.shift()
		return float, nil
	case StringTok:
		str := newString(string(emitter.curr().buf))
		emitter.shift()
		return str, nil
	case LParenTok:
		emitter.shift()
		var elements []ast
		for emitter.curr().typ != RParenTok {
			elem, err := parse(emitter)
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)
		}
		emitter.shift()
		return newList(elements), nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", emitter.curr())
	}
}