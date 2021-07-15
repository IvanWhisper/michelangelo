package excode

import "testing"

func TestCode(t *testing.T) {
	for i := 0; i < 5; i++ {
		println(ExCode(i).ToEx("123").ToError().Error())
	}
}
