package fileex

import (
	"testing"
)

func TestPathExistOrNot(t *testing.T) {
	p := "./test.txt"
	b, err := PathExistOrNot(p)
	if err != nil {
		t.Error(err)
	}
	if b {
		t.Error(p, " is exist")
	}
	t.Log(p, " is not exist")
}

func ExamplePathExistOrNot() {
	fPath := "./test.txt"
	b, err := PathExistOrNot(fPath)
	if err != nil {
		println(err)
	}
	if b {
		println("file path is exist")
	}
}

func BenchmarkPathExistOrNot(b *testing.B) {

}
