package main

import (
	"errors"
	"fmt"
	"io"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"

	"github.com/alesiong/codec/codecs"
)

type charsetCodecs struct {
}

func (c charsetCodecs) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	charset := options["C"]
	if charset == "" {
		return errors.New("charset: missing required option charset (-C)")
	}

	// toCharset := options["T"]

	encoders := map[string]encoding.Encoding{
		"gbk":       simplifiedchinese.GBK,
		"big5":      traditionalchinese.Big5,
		"shift-jis": japanese.ShiftJIS,
	}

	encoder, ok := encoders[charset]
	if !ok {
		return fmt.Errorf("charset: unkonwn charset %s", charset)
	}

	switch globalMode {
	case codecs.CodecModeEncoding:
		err = codecs.ReadToWriter(input, transform.NewWriter(output, encoder.NewEncoder()))
	case codecs.CodecModeDecoding:
		err = codecs.ReadToWriter(transform.NewReader(input, encoder.NewDecoder()), output)
	default:
		return errors.New("invalid codec mode")
	}
	return
}

var CodecPlugin codecs.Codec = charsetCodecs{}

func main() {
}
