package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alesiong/codec/codecs"
)

const (
	optionEncoding    = "e"
	optionDecoding    = "d"
	optionNewLine     = "n"
	optionInputString = "I"
	optionInputFile   = "F"
	optionOutputFile  = "O"
	optionHelp        = "h"
)

type executor struct{}

func (e executor) execute(command *command) (err error) {
	globalMode := codecs.CodecModeEncoding

loop:
	for _, o := range command.options {
		switch o.name {
		case optionDecoding:
			globalMode = codecs.CodecModeDecoding
		case optionNewLine:
			newLineCodec := codec{
				name: "newline",
			}
			command.codecs = append(command.codecs, newLineCodec)
		case optionInputString:
			constCodec := codec{
				name: "const",
				options: []option{
					{
						optionType: optionTypeValue,
						name:       "C",
						text:       o.text,
					},
				},
			}
			command.codecs = append([]codec{constCodec}, command.codecs...)
		case optionInputFile:
			catCodec := codec{
				name: "cat",
				options: []option{
					{
						optionType: optionTypeSwitch,
						name:       "c",
					},
					{
						optionType: optionTypeValue,
						name:       "F",
						text:       o.text,
					},
				},
			}
			command.codecs = append([]codec{catCodec}, command.codecs...)
		case optionOutputFile:
			sinkCodec := codec{
				name: "redirect",
				options: []option{
					{
						optionType: optionTypeValue,
						name:       "O",
						text:       o.text,
					},
				},
			}
			command.codecs = append(command.codecs, sinkCodec)
		case optionHelp:
			usageCodec := codec{
				name: "usage",
			}
			command.codecs = []codec{usageCodec}
			break loop
		default:
			return fmt.Errorf("unknown option: %s", o.name)
		}

	}

	err = e.runCodecs(os.Stdin, command.codecs, globalMode, os.Stdout)
	if err != nil {
		return
	}

	return
}

func (e executor) runCodecs(input io.Reader, codecList []codec, mode codecs.CodecMode, output io.Writer) (err error) {
	previousInput := input
	for _, c := range codecList {
		reader, writer := io.Pipe()
		go func(input io.Reader, output io.WriteCloser, c codec) {
			defer output.Close()
			err = e.runCodec(input, &c, mode, output)
			if err != nil {
				err = fmt.Errorf("error in %s: %w", c.name, err)
				panic(err)
			}
		}(previousInput, writer, c)

		previousInput = reader
	}

	_, err = io.Copy(output, previousInput)
	return
}

func (e executor) runCodec(input io.Reader, codec *codec, mode codecs.CodecMode, output io.Writer) (err error) {
	options, err := e.makeCodecOptions(codec)
	if err != nil {
		return
	}
	if options["e"] != "" {
		mode = codecs.CodecModeEncoding
	}
	if options["d"] != "" {
		mode = codecs.CodecModeDecoding
	}

	if c := codecs.Lookup(codec.name); c != nil {
		return c.RunCodec(input, mode, options, output)
	} else {
		return fmt.Errorf("codec not found: %s", codec.name)
	}
}

func (e executor) makeCodecOptions(codec *codec) (option map[string]string, err error) {
	option = make(map[string]string)

	for _, o := range codec.options {
		switch o.optionType {
		case optionTypeSwitch:
			option[o.name] = "*" // TODO: eliminate hard coding
		case optionTypeValue:
			switch o.text.textType {
			case textTypeString:
				option[o.name] = o.text.string
			case textTypeCodec:
				var buf bytes.Buffer
				err := e.runCodecs(strings.NewReader(o.text.string), o.text.codecs, codecs.CodecModeEncoding, &buf)
				if err != nil {
					return nil, err
				}
				output, err := ioutil.ReadAll(&buf)
				if err != nil {
					return nil, err
				}
				option[o.name] = string(output)
			}
		}
	}
	return
}
