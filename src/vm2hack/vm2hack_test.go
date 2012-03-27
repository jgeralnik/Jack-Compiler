package vm2hack_test

import (
	"bytes"
	"testing"
	. "vm2hack"
)

func TestConvertLine(t *testing.T) {
	result_ := ConvertLine([]byte("push constant 5"))
	if !bytes.Equal(result, []byte("@5\nM=A\n...")) {
		t.Errorf("Failed translating push constant 5")
	}
}

func TestConvertFile(t *testing.T) {

}
