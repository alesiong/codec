package builtins

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"strings"

	"github.com/alesiong/codec/codecs"
)

type aesCodec struct {
	mode string
}

func init() {
	codecs.Register("aes-cbc", aesCodec{mode: "cbc"})
	codecs.Register("aes-ecb", aesCodec{mode: "ecb"})
}

func (b aesCodec) Usage() string {
	switch b.mode {
	case "cbc":
		return strings.TrimLeft(`
    -K key
    -IV iv
`, "\n")
	case "ecb":
		return strings.TrimLeft(`
    -K key
`, "\n")
	}
	return ""
}

func (b aesCodec) RunCodec(input io.Reader, globalMode codecs.CodecMode, options map[string]string, output io.Writer) (err error) {
	key := options["K"]
	if key == "" {
		return errors.New("aes: missing required option key (-K)")
	}

	crypto, err := aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	blockSize := crypto.BlockSize()
	encryptFunc := crypto.Encrypt
	decryptFunc := crypto.Decrypt

	if b.mode == "cbc" {
		iv := options["IV"]
		if iv == "" {
			return errors.New("aes[cbc]: missing required option iv (-IV)")
		}
		var c cipher.BlockMode
		switch globalMode {
		case codecs.CodecModeEncoding:
			c = cipher.NewCBCEncrypter(crypto, []byte(iv))
		case codecs.CodecModeDecoding:
			c = cipher.NewCBCDecrypter(crypto, []byte(iv))
		default:
			return errors.New("invalid codec mode")
		}

		blockSize = c.BlockSize()
		encryptFunc = c.CryptBlocks
		decryptFunc = c.CryptBlocks
	}

	switch globalMode {
	case codecs.CodecModeEncoding:
		err = encrypt(input, output, encryptFunc, blockSize, pkcs5Pad)
	case codecs.CodecModeDecoding:
		err = decrypt(input, output, decryptFunc, blockSize, pkcs5Unpad)
	default:
		return errors.New("invalid codec mode")
	}
	return
}

func pkcs5Pad(src []byte, blockSize int) []byte {
	padding := aes.BlockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs5Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[len(src)-1])
	for _, b := range src[length-unpadding:] {
		if int(b) != unpadding {
			return nil, errors.New("pkcs: unpadding failed")
		}
	}
	return src[:length-unpadding], nil
}

func encrypt(input io.Reader, output io.Writer, encryptFunc func([]byte, []byte), blockSize int, pad func([]byte, int) []byte) (err error) {
	block := make([]byte, blockSize)
	for {
		n, e := io.ReadFull(input, block)
		if e == io.EOF || e == io.ErrUnexpectedEOF {
			block = pad(block[:n], blockSize)
			break
		}

		if e != nil {
			return e
		}

		encryptFunc(block, block)
		_, err = output.Write(block)
	}

	encryptFunc(block, block)
	_, err = output.Write(block)
	return
}

func decrypt(input io.Reader, output io.Writer, decryptFunc func([]byte, []byte), blockSize int, unpad func([]byte) ([]byte, error)) (err error) {
	var lastBlock []byte
	block := make([]byte, blockSize)

	for {
		_, e := io.ReadFull(input, block)
		if e == io.EOF {
			break
		}
		if lastBlock == nil {
			lastBlock = make([]byte, blockSize)
		} else {
			_, err = output.Write(lastBlock)
			if err != nil {
				return
			}
		}

		if e != nil {
			return e
		}
		decryptFunc(lastBlock, block)
	}

	lastBlock, err = unpad(lastBlock)
	if err != nil {
		return
	}

	_, err = output.Write(lastBlock)
	return
}
