package codecs

import (
	"io"
)

type CodecMode int

const (
	CodecModeEncoding CodecMode = iota
	CodecModeDecoding
)

type Codec interface {
	Execute(input io.Reader, globalMode CodecMode, options map[string]string, output io.WriteCloser) (err error)
}

var (
	AesCbc Codec = aesCodec{mode: "cbc"}
	AesEcb Codec = aesCodec{mode: "ecb"}
	Base64 Codec = base64Codec{}
	Url    Codec = urlCodecs{}
	Zlib   Codec = zlibCodec{}
	Hex    Codec = hexCodecs{}
	Md5    Codec = hashCodecs{mode: "md5"}
	Sha256 Codec = hashCodecs{mode: "sha256"}
	Escape Codec = escapeCodecs{}

	Id       Codec = idCodecs{}
	Const    Codec = constCodecs{}
	Repeat   Codec = repeatCodecs{}
	Tee      Codec = teeCodecs{}
	Redirect Codec = redirectCodecs{}
	Sink     Codec = sinkCodecs{}
	Append   Codec = appendCodecs{}
	Newline  Codec = newLineCodecs{}
)
