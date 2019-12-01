package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"strings"

	_ "github.com/alesiong/codec/codecs/builtins"
	_ "github.com/alesiong/codec/codecs/meta"
)

func main() {
	tokenizer := tokenizer{
		text: os.Args[1:],
	}

	command := parseCommand(&tokenizer)

	loadPlugins()

	executor := executor{}

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

func loadPlugins() {
	files, err := ioutil.ReadDir("plugins")
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "error when loading plugins:", err)
		}
		return
	}
	for _, f := range files {
		const pluginExtension = ".so"
		if path.Ext(f.Name()) != pluginExtension {
			continue
		}
		pluginKey := strings.TrimSuffix(f.Name(), pluginExtension)
		_, err := plugin.Open("plugins/" + f.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error when loading plugin %s: %v\n", pluginKey, err)
			continue
		}
	}

}
