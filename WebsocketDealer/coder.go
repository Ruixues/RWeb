package WebsocketDealer

import (
	"bytes"

	"github.com/klauspost/compress/gzip"
)

type Coder interface {
	Encode([]byte) []byte
	Decode([]byte) ([]byte, error)
}
type DefaultCoder struct{}

func (z *DefaultCoder) Encode(data []byte) []byte {
	var compressData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressData)
	defer gzipWriter.Close()
	gzipWriter.Write(data)
	gzipWriter.Flush()
	return compressData.Bytes()
}
func Decode(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	var ret []byte
	if _, err := reader.Read(ret); err != nil {
		return nil, err
	}
	return ret, err
}
