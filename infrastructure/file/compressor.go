package fileex

import "os"

type CompressType int32

const (
	CompressNil CompressType = iota
	CompressZip
)

type Compressor interface {
	Compress(output string, files ...*os.File) error
	UnCompress(zipFile, output string) error
}

// NewCompressor
/**
 * @Description:
 * @param compressType
 * @return Compressor
 */
func NewCompressor(compressType CompressType) Compressor {
	switch compressType {
	case CompressZip:
		return &ZipCompress{
			Type: compressType,
		}
	case CompressNil:
		return nil
	default:
		return nil
	}
}
