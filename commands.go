package main

import "fmt"

type command struct {
	options []option
	codecs  []codec
	string  string
}

func (c *command) print() {
	fmt.Print("Options: ")
	for _, option := range c.options {
		option.print()
		fmt.Print("\t")
	}
	fmt.Println()

	fmt.Print("Codecs: ")
	for _, codec := range c.codecs {
		codec.print()
		fmt.Print("\t")
	}
	fmt.Println()

	// fmt.Println("String: ", c.string)
}

type codec struct {
	name    string
	options []option
}

func (c *codec) print() {
	fmt.Print(c.name)
	fmt.Print("(")
	for _, option := range c.options {
		option.print()
		fmt.Print("\t")
	}
	fmt.Print(")")
}

type textType int

const (
	textTypeString textType = iota
	textTypeCodec
)

type text struct {
	textType textType
	string   string
	codecs   []codec // valid if textTypeCodec
}

func (t *text) print() {
	switch t.textType {
	case textTypeString:
		fmt.Print(t.string)
	case textTypeCodec:
		fmt.Print("[ ")
		for _, codec := range t.codecs {
			codec.print()
			fmt.Print("\t")
		}
		fmt.Print(t.string)
		fmt.Print(" ]")
	}
}

type optionType int

const (
	optionTypeSwitch optionType = iota
	optionTypeValue
)

type option struct {
	optionType optionType
	name       string
	text       *text // valid if optionTypeValue
}

func (o *option) print() {
	switch o.optionType {
	case optionTypeValue:
		fmt.Print(o.name, " ")
		o.text.print()
	case optionTypeSwitch:
		fmt.Print(o.name)
	}
}
