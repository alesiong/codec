package main

import (
	"bytes"
	"errors"
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
)

type executor struct {
	codecsMap map[string]codecs.Codec
}

func (e *executor) execute(command *command) (err error) {
	globalMode := codecs.CodecModeEncoding

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
			if o.text.textType != textTypeString {
				return errors.New("main option -I cannot have codecs syntax")
			}
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

		case optionOutputFile:
			if o.text.textType != textTypeString {
				return errors.New("main option -O cannot have codecs syntax")
			}
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

func (e *executor) runCodecs(input io.Reader, codecList []codec, mode codecs.CodecMode, output io.Writer) (err error) {
	previousInput := input
	buffers := make([]codecs.BlockingReader, 0, len(codecList))
	for _, c := range codecList {
		buf := codecs.BlockingReader{Chan: make(chan []byte)}
		go func(input io.Reader, output io.WriteCloser, c codec) {
			defer output.Close()
			err = e.runCodec(input, &c, mode, output)
			if err != nil {
				panic(err)
			}
		}(previousInput, &buf, c)

		previousInput = &buf
		buffers = append(buffers, buf)
	}

	err = codecs.ReadToWriter(previousInput, output)
	return
}

func (e *executor) runCodec(input io.Reader, codec *codec, mode codecs.CodecMode, output io.Writer) (err error) {
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

	if c, ok := e.codecsMap[codec.name]; ok {
		return c.RunCodec(input, mode, options, output)
	} else {
		return fmt.Errorf("codec not found: %s", codec.name)
	}
}

func (e *executor) makeCodecOptions(codec *codec) (option map[string]string, err error) {
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
