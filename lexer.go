package main

import (
	"fmt"
	"regexp"
)

//go:generate stringer -type=tokenType

type tokenType int

const (
	LParenTok tokenType = iota
	RParenTok
	AtomTok
	IntTok
	FloatTok
	StringTok

	QuoteTok
	QuasiQuoteTok
	UnquoteTok
)

type tokenizer struct {
	buf []byte
	pos int
}

func newTokenizer(buf []byte) *tokenizer {
	return &tokenizer{buf: buf}
}

func (t *tokenizer) eof() bool {
	return t.pos >= len(t.buf)
}

func (t *tokenizer) match(i *regexp.Regexp) ([]byte, bool) {
	loc := i.FindIndex(t.buf[t.pos:])
	if loc == nil {
		return nil, false
	}
	if loc[0] != 0 {
		return nil, false
	}
	return t.buf[t.pos:][loc[0]: loc[1]], true
}

func (t *tokenizer) advance(n int) {
	t.pos += n
}

type token struct {
	typ tokenType
	buf []byte
}

func (t token) String() string {
	return fmt.Sprintf("[%s: %q]", t.typ, t.buf)
}

func lex(source *tokenizer) ([]token, error) {
	var recognizedTokens = []struct{tokenType; *regexp.Regexp}{
		{LParenTok, regexp.MustCompile(`\(`)},
		{RParenTok, regexp.MustCompile(`\)`)},

		{QuoteTok, regexp.MustCompile(`'`)},
		{QuasiQuoteTok, regexp.MustCompile("`")},
		{UnquoteTok, regexp.MustCompile(`,`)},

		{FloatTok,  regexp.MustCompile(`\d+\.\d+`)},
		{IntTok,    regexp.MustCompile(`\d+`)},
		{StringTok, regexp.MustCompile(`"[^"]*"`)},
		{AtomTok,   regexp.MustCompile(`[a-zA-Z=+\-/]+`)},
	}

	whitespace := regexp.MustCompile(`^\s`)

	var tokens []token
outer:
	for !source.eof() {
		if buf, ok := source.match(whitespace); ok {
			source.advance(len(buf))
			continue outer
		}

		for _, option := range recognizedTokens {
			if buf, ok := source.match(option.Regexp); ok {
				tokens = append(tokens, token{
					typ: tokenType(option.tokenType),
					buf: buf,
				})
				source.advance(len(buf))
				continue outer
			}
		}
		return nil, fmt.Errorf("unexpected token at position: %d", source.pos)
	}
	return tokens, nil
}
