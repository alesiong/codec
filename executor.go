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
	optionEncoding = "e"
	optionDecoding = "d"
	optionNewLine  = "n"
)

type executor struct {
	codecsMap map[string]codecs.Codec
}

func (e *executor) execute(command *command, output io.Writer) (err error) {
	globalMode := codecs.CodecModeEncoding
	endWithNewLine := false

	for _, option := range command.options {
		switch option.name {
		case optionDecoding:
			globalMode = codecs.CodecModeDecoding
		case optionNewLine:
			endWithNewLine = true
		}

	}

	var input io.Reader

	if command.string == "" {
		input = os.Stdin
	} else {
		input = strings.NewReader(command.string)
	}

	err = e.runCodecs(input, command.codecs, globalMode, output)
	if err != nil {
		return
	}

	if endWithNewLine {
		_, err = output.Write([]byte{'\n'})
	}
	return
}

func (e *executor) runCodecs(input io.Reader, codecList []codec, mode codecs.CodecMode, output io.Writer) (err error) {
	previousInput := input
	buffers := make([]codecs.BlockingReader, 0, len(codecList))
	for _, c := range codecList {
		buf := codecs.BlockingReader{Chan: make(chan []byte)}
		go func(input io.Reader, output io.WriteCloser, c codec) {
			err = e.runCodec(input, &c, mode, output)
			if err != nil {
				panic(err)
			}
		}(previousInput, &buf, c)

		previousInput = &buf
		buffers = append(buffers, buf)
	}

	err = codecs.ReadToWriter(previousInput, output, nil)
	return
}

func (e *executor) runCodec(input io.Reader, codec *codec, mode codecs.CodecMode, output io.WriteCloser) (err error) {
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
		return c.Execute(input, mode, options, output)
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
