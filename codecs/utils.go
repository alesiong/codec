package codecs

import (
	"io"
)

const (
	bufferSize = 1
)

func ReadToWriter(reader io.Reader, writer io.Writer, closer io.Closer) (err error) {
	buf := make([]byte, bufferSize)
	n := 0
	for {
		n, err = reader.Read(buf)
		if err != nil && err != io.EOF {
			return
		}

		data := buf[:n]

		// for _, t := range transformer {
		// 	var e error
		// 	data, e = t(data)
		// 	if e != nil {
		// 		return e
		// 	}
		// }

		_, e := writer.Write(data)
		if e != nil {
			return e
		}

		if err == io.EOF {
			break
		}
	}
	if closer != nil {
		return closer.Close()
	}
	return nil
}

type BlockingReader struct {
	Chan   chan []byte
	buffer []byte
}

func (b *BlockingReader) Write(p []byte) (n int, err error) {
	length := len(p)
	buf := make([]byte, length)
	copy(buf, p)
	b.Chan <- buf
	return length, nil
}

func (b *BlockingReader) Read(p []byte) (n int, err error) {
	var buf []byte
	if len(b.buffer) == 0 {
		ok := false
		buf, ok = <-b.Chan
		if !ok {
			return 0, io.EOF
		}
	} else {
		buf = b.buffer
	}

	n = copy(p, buf)
	if n < len(buf) {
		b.buffer = buf[n:]
	} else {
		b.buffer = nil
	}
	return
}

func (b *BlockingReader) Close() error {
	close(b.Chan)
	return nil
}
