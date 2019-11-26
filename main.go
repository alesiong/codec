package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"strings"

	"github.com/alesiong/codec/codecs"
	"github.com/alesiong/codec/codecs/builtins"
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
	builtIns := map[string]codecs.Codec{
		"aes-cbc":  builtins.AesCbc,
		"aes-ecb":  builtins.AesEcb,
		"base64":   builtins.Base64,
		"hex":      builtins.Hex,
		"url":      builtins.Url,
		"sha256":   builtins.Sha256,
		"md5":      builtins.Md5,
		"zlib":     builtins.Zlib,
		"id":       builtins.Id,
		"const":    builtins.Const,
		"repeat":   builtins.Repeat,
		"tee":      builtins.Tee,
		"redirect": builtins.Redirect,
		"sink":     builtins.Sink,
		"append":   builtins.Append,
		"newline":  builtins.Newline,
		"escape":   builtins.Escape,
		"cat":      builtins.Cat,
	}
	loadPlugins(builtIns)
	return builtIns
}

func loadPlugins(codecsMap map[string]codecs.Codec) {
	files, err := ioutil.ReadDir("plugins")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error when loading plugins:", err)
	}
	for _, f := range files {
		const pluginExtension = ".so"
		if path.Ext(f.Name()) != pluginExtension {
			continue
		}
		pluginKey := strings.TrimSuffix(f.Name(), pluginExtension)
		p, err := plugin.Open("plugins/" + f.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error when loading plugin %s: %v\n", pluginKey, err)
			continue
		}
		codec, err := p.Lookup("CodecPlugin")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error when loading plugin %s: %v\n", pluginKey, err)
			continue
		}
		c, ok := codec.(*codecs.Codec)
		if ok {
			codecsMap[pluginKey] = *c
		} else {
			fmt.Fprintf(os.Stderr, "error when loading plugin %s: %v\n", pluginKey, err)
			continue
		}
	}

}
