package exec

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	req := ExecReq{}
	out, cancel := Start(&req)
	defer cancel()
	for v := range out {
		fmt.Println(v)
	}
}
