package main

import (
	"os"

	"github.com/alesiong/codec/codecs"
)

func main() {
	tokenizer := tokenizer{
		text: os.Args[1:],
	}

	command := parseCommand(&tokenizer)

	codecsMap := loadCodecs()

	executor := executor{
		codecsMap: codecsMap,
	}

	// TODO: eliminate panic
	err := executor.execute(&command)
	if err != nil {
		panic(err)
	}
}

func parseCommand(tokenizer *tokenizer) (command command) {
	options, err := parseOptions(tokenizer)
	if err != nil {
		panic(err)
	}
	command.options = options
	for {
		codec, err := parseCodec(tokenizer)
		if err != nil {
			panic(err)
		}

		if codec == nil {
			break
		}

		command.codecs = append(command.codecs, *codec)
	}
	return
}

func loadCodecs() map[string]codecs.Codec {
	return map[string]codecs.Codec{
		"aes-cbc":  codecs.AesCbc,
		"aes-ecb":  codecs.AesEcb,
		"base64":   codecs.Base64,
		"hex":      codecs.Hex,
		"url":      codecs.Url,
		"sha256":   codecs.Sha256,
		"md5":      codecs.Md5,
		"zlib":     codecs.Zlib,
		"id":       codecs.Id,
		"const":    codecs.Const,
		"repeat":   codecs.Repeat,
		"tee":      codecs.Tee,
		"redirect": codecs.Redirect,
		"sink":     codecs.Sink,
		"append":   codecs.Append,
		"newline":  codecs.Newline,
		"escape":   codecs.Escape,
	}
}
