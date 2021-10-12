package fileex

import (
	"path/filepath"
	"testing"
)

func TestMkdirIfNotExist(t *testing.T) {
	p := "./test.txt"
	absPath, _ := filepath.Abs(p)
	t.Log("temp file is", absPath)
	defer PathDelete(p)
	err := MkdirIfNotExist(p)
	if err != nil {
		t.Error(err)
	}
	b, err := PathExistOrNot(p)
	if err != nil {
		t.Error(err)
	}
	if !b {
		t.Error(p, " create failed")
	}
	t.Log(p, " create success")
}
