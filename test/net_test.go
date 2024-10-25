package test

import (
	"testing"

	"github.com/kbrownehs18/goutil/net"
)

func TestMask(t *testing.T) {
	ip, err := net.BitsToMask(25)
	t.Error(err)
	t.Logf("IP: %s", ip)

	bit, err := net.MaskToBits("255.255.0.0")
	t.Error(err)
	t.Logf("Bit: %d", bit)
}
