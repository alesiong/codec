package meta

import (
	"fmt"
	"io"

	"github.com/alesiong/codec/codecs"
)

type usageMetaCodec struct {
}

func init() {
	codecs.RegisterMeta("usage", usageMetaCodec{})
}

func (u usageMetaCodec) RunMetaCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, metaInfo *codecs.CodecMetaInfo, output io.Writer) (err error) {
	if _, err := fmt.Fprintln(output, "Available codecs:"); err != nil {
		return err
	}
	nameCh, codecCh := metaInfo.AllCodecs()
	for {
		name, ok := <-nameCh
		if !ok {
			break
		}

		codec := <-codecCh

		if _, err := fmt.Fprintln(output, name); err != nil {
			return err
		}

		if usage, ok := codec.(codecs.CodecUsage); ok {
			if _, err := fmt.Fprintln(output, usage.Usage()); err != nil {
				return err
			}
		}
	}
	return
}
