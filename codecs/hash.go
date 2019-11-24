package codecs

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"hash"
	"io"
)

type hashCodecs struct {
	mode string
}

func (h hashCodecs) RunCodec(input io.Reader, globalMode CodecMode, options map[string]string, output io.Writer) (err error) {
	var hasher hash.Hash
	switch h.mode {
	case "md5":
		hasher = md5.New()
	case "sha256":
		hasher = sha256.New()

	default:
		return errors.New("hash: invalid mode")
	}

	switch globalMode {
	case CodecModeEncoding:
		err = ReadToWriter(input, hasher)
		if err != nil {
			return
		}
		_, err = output.Write(hasher.Sum(nil))
	case CodecModeDecoding:
		return errors.New("hash: cannot decode")
	default:
		return errors.New("invalid codec mode")
	}

	return
}
