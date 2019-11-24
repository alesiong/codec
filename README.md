# codec
A Swiss Army knife for en/de- coding/crypting strings

## Description
`codec` will read string from `stdin`, transform it through multiple en/de-coders
and print the result to `stdout`

## Usage
```
codec options codecs

options:
    -e (default): set the global coding mode to encode
    -d: set the global coding mode to decode
    -n: append new line ('\n') at the end of output

codecs:
    a list of codecs(en/de-coders), input will be passed and transformed from
    left to right

codec:
    codec-name codec-options

codec-options:
    lower case options are switch(boolean) options, so they take no argument.

    upper case options take one argument. the argument can be provided with plain
    string or by sub-codecs syntax: [codecs plain-string]. If you use sub-codecs
    syntax, the codecs inside [] will be run on `plain-string` as input, and the
    output is used as the argument.
```

### Examples
You can use `echo -n '' | ` to pass the input string directly.
```
codec -d base64 zlib
```
Decode base64 on input, and then decompress with zlib.

```
codec aes-ecb -K [hex -d 12345678901234561234567890123456] base64
```
Decode hex string `12345678901234561234567890123456`, and set it as aes-ecb key.
Encrypt the input, and then encode with base64. Note that unlike `openssl`, aes
codecs do not expect hex string as key. You always pass a raw byte string as key.

### Available Codecs and Options
If `-d` or `-e` is passed as a codec option, it will overwrite the global coding
mode.

```
url
    url query escape/unescape
    -p: use path escape instead of query escape

base64
    -u: use url base64 instead

aes-ecb
    -K key

aes-cbc
    -K key
    -IV iv

hex
    binary to hex encode or inverse
    -c: use capital hex string (only affects encoding)

sha256

md5

zlib
    -L level: compress level (int, [-2, 9])
```

# TODO
1. refactor code
2. load codecs as go plugins