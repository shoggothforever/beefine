package image_prep

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	req := ImagePrepReq{}
	out, cancel := Start(&req)
	defer cancel()
	for v := range out {
		fmt.Println(v)
	}
}
