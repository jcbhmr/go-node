package punycode_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/jcbhmr/go-node/punycode"
)

func try1[A any](a A, err error) A {
	if err != nil {
		panic(err)
	}
	return a
}

func execNodePrint(code string) string {
	cmd := exec.Command("node", "--print", code)
	bytes := try1(cmd.CombinedOutput())
	return strings.TrimSuffix(string(bytes), "\n")
}

func TestToAscii(t *testing.T) {
	scenarios := map[string]string{
		"hello":         "",
		"world":         "",
		"📗":             "",
		"👈🛑❤🪀🚀🥳":        "",
		"🙌.example.org": "",
		"🧰@example.org": "",
	}
	for input, expected := range scenarios {
		js := fmt.Sprintf("punycode.toASCII(%#v)", input)
		expected = execNodePrint(js)
		scenarios[input] = expected
	}
	for input, expected := range scenarios {
		actual := punycode.ToAscii(input)
		if actual != expected {
			t.Errorf("expected %#v, got %#v", expected, actual)
		}
	}
}
