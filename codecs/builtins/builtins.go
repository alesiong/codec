package builtins

import "github.com/alesiong/codec/codecs"

var (
	AesCbc codecs.Codec = aesCodec{mode: "cbc"}
	AesEcb codecs.Codec = aesCodec{mode: "ecb"}
	Base64 codecs.Codec = base64Codec{}
	Url    codecs.Codec = urlCodecs{}
	Zlib   codecs.Codec = zlibCodec{}
	Hex    codecs.Codec = hexCodecs{}
	Md5    codecs.Codec = hashCodecs{mode: "md5"}
	Sha256 codecs.Codec = hashCodecs{mode: "sha256"}
	Escape codecs.Codec = escapeCodecs{}

	Id       codecs.Codec = idCodecs{}
	Const    codecs.Codec = constCodecs{}
	Repeat   codecs.Codec = repeatCodecs{}
	Tee      codecs.Codec = teeCodecs{}
	Redirect codecs.Codec = redirectCodecs{}
	Sink     codecs.Codec = sinkCodecs{}
	Append   codecs.Codec = appendCodecs{}
	Newline  codecs.Codec = newLineCodecs{}
	Cat      codecs.Codec = catCodecs{}
)
