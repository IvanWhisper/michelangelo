package fileex

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

// MockTempFiles
/**
 * @Description:
 * @param exec
 */
func MockTempFiles(exec func(paths ...string)) {
	defer func() {
		err := PathDelete("./testTemp")
		if err != nil {
			panic(err)
		}
	}()
	tmpPath := path.Join("./testTemp", "tmp")
	err := MkdirIfNotExist(tmpPath)
	if err != nil {
		panic(err)
	}
	defer PathDelete(tmpPath)
	file1 := path.Join("./testTemp", "1.txt")
	err = ioutil.WriteFile(file1, []byte("123456A"), fs.ModePerm)
	if err != nil {
		panic(err)
	}
	defer PathDelete(file1)
	file2 := path.Join(tmpPath, "2.txt")
	err = ioutil.WriteFile(file2, []byte("123456B"), fs.ModePerm)
	if err != nil {
		panic(err)
	}
	defer PathDelete(file2)
	exec(tmpPath, file1)
}

// TestCompress
/**
 * @Description:
 * @param t
 */
func TestZipCompress(t *testing.T) {
	compressor := NewCompressor(CompressZip)
	MockTempFiles(func(paths ...string) {
		output := path.Join("./testTemp", "tmp.zip")
		defer PathDelete(output)
		files := make([]*os.File, 0)
		for _, p := range paths {
			f, err := os.Open(p)
			if err != nil {
				t.Fatal(err)
			}
			files = append(files, f)
		}
		defer func() {
			for _, f := range files {
				f.Close()
			}
		}()
		beg := time.Now()
		err := compressor.Compress(output, files...)
		t.Logf("compress cost %dms", time.Since(beg).Milliseconds())
		if err != nil {
			t.Fatal(err)
		}
		beg2 := time.Now()
		err = compressor.UnCompress(output, "./testTemp/tmp")
		t.Logf("uncompress cost %dms", time.Since(beg2).Milliseconds())
		if err != nil {
			t.Fatal(err)
		}
	})
}
