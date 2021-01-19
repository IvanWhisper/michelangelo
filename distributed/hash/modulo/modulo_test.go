package modulo

import (
	"strconv"
	"testing"
)

func Test_GetShardIndex(t *testing.T) {
	s := New(3)
	index := s.GetShardIndex(10)
	if index != 1 {
		t.Errorf("10 GetShardIndex => %s!=1", strconv.FormatUint(index, 10))
	} else {
		t.Log("GetShardIndex Success")
	}

}
