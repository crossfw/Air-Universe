package SysLoad

import (
	"fmt"
	"testing"
)

func TestGetSysLoad(t *testing.T) {
	a, err := GetSysLoad()
	if err != nil {
		t.Errorf("Failed")
	}
	fmt.Println(a)
}
