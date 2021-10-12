package fileex

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

var _ = Compressor(&ZipCompress{})

// ZipCompress "Zip Compress"
type ZipCompress struct {
	Type CompressType
}

// Compress
/**
 * @Description:
 * @receiver c
 * @param output
 * @param files
 * @return error
 */
func (c *ZipCompress) Compress(output string, files ...*os.File) error {
	compressionFile, err := os.Create(output)
	defer compressionFile.Close() // nolint:gosec
	if err != nil {
		return err
	}
	w := zip.NewWriter(compressionFile)
	defer w.Close()
	for _, file := range files {
		err := c.compress(file, w)
		if err != nil {
			return err
		}
	}
	return nil
}

// compress
/**
 * @Description:
 * @receiver c
 * @param file
 * @param zw
 * @return error
 */
func (c *ZipCompress) compress(file *os.File, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(path.Join(file.Name(), fi.Name()))
			if err != nil {
				return err
			}
			err = c.compress(f, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// UnCompress
/**
 * @Description:
 * @receiver c
 * @param zipFile
 * @param output
 * @return error
 */
func (c *ZipCompress) UnCompress(zipFile, output string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		err := c.unCompress(file, output)
		if err != nil {
			return err
		}
	}
	return nil
}

// unCompress
/**
 * @Description:
 * @receiver c
 * @param file
 * @param output
 * @return error
 */
func (c *ZipCompress) unCompress(file *zip.File, output string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	err = MkdirIfNotExist(output)
	if err != nil {
		return err
	}
	w, err := os.Create(path.Join(output, file.Name)) //nolint:gosec
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, rc) //nolint:gosec
	if err != nil {
		return err
	}
	return nil
}
