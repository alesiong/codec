package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	openingParenthesis = "["
	closingParenthesis = "]"

	optionPrefix = "-"
)

type tokenizer struct {
	lookNext   string
	text       []string
	currentPos int
	eof        bool
}

func (p *tokenizer) next() (next string) {
	if p.lookNext != "" {
		next = p.lookNext
		p.lookNext = ""
		return
	}

	if p.eof {
		return
	}

	next = p.text[p.currentPos]
	p.currentPos++

	if p.currentPos >= len(p.text) {
		p.eof = true
	}

	if strings.HasPrefix(next, openingParenthesis) && next != openingParenthesis {
		p.lookNext = next[len(openingParenthesis):]
		return openingParenthesis
	}

	if strings.HasSuffix(next, closingParenthesis) && next != closingParenthesis {
		p.lookNext = closingParenthesis
		return next[:len(next)-len(closingParenthesis)]
	}

	return
}

func (p *tokenizer) peek() string {
	if p.lookNext != "" {
		return p.lookNext
	}

	if p.currentPos >= len(p.text) {
		p.eof = true
	}

	if p.eof {
		return ""
	}

	next := p.text[p.currentPos]
	if strings.HasPrefix(next, openingParenthesis) && next != openingParenthesis {
		return openingParenthesis
	}

	if strings.HasSuffix(next, closingParenthesis) && next != closingParenthesis {
		return next[:len(next)-len(closingParenthesis)]
	}
	return next
}

func isSpecialToken(token string) bool {
	return token == openingParenthesis || token == closingParenthesis
}

func parseOptions(tokenizer *tokenizer) (options []option, err error) {
	for {
		optionName := tokenizer.peek()
		if !strings.HasPrefix(optionName, optionPrefix) {
			break
		}
		tokenizer.next()

		var option option
		option.name = optionName[len(optionPrefix):]

		firstRune, _ := utf8.DecodeRuneInString(option.name)
		if unicode.IsUpper(firstRune) {
			// value option
			option.optionType = optionTypeValue

			var text *text
			text, err = parseText(tokenizer)
			if err != nil {
				options = nil
				return
			}

			option.text = text
		} else {
			// switch option
			option.optionType = optionTypeSwitch
		}

		options = append(options, option)
	}
	return
}

func parseCodec(tokenizer *tokenizer) (codec_ *codec, err error) {
	name := tokenizer.peek()
	if name == "" || isSpecialToken(name) {
		return
	}

	tokenizer.next()

	options, err := parseOptions(tokenizer)
	if err != nil {
		return
	}

	codec_ = &codec{
		name:    name,
		options: options,
	}
	return
}

func parseText(tokenizer *tokenizer) (text_ *text, err error) {
	str := tokenizer.peek()

	// TODO: allow empty option value
	if str == "" {
		err = errors.New("EOF when parsing")
		return
	}

	tokenizer.next()

	text_ = new(text)

	// TODO: escape parenthesis
	// TODO: allow empty sub-codecs syntax
	if str == openingParenthesis {
		text_.textType = textTypeCodec
		text_.string = tokenizer.next()

		for {
			codec, err := parseCodec(tokenizer)
			if err != nil {
				return nil, err
			}

			if codec == nil {
				break
			}

			text_.codecs = append(text_.codecs, *codec)
		}
		if n := tokenizer.next(); n != closingParenthesis {
			return nil, fmt.Errorf("expect %s, found %s", closingParenthesis, n)
		}
	} else {
		text_.textType = textTypeString
		text_.string = str
	}

	return
}
