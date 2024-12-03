package counter

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	req := CounterReq{IfName: "ens33"}
	out, cancel := Start(&req)
	defer cancel()
	for v := range out {
		fmt.Println(v)
	}
}
