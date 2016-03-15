package beater

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseMockOutput(output string) chan pkg {
	return parseOutput(ioutil.NopCloser(strings.NewReader(output)))
}

func TestParseOutput(t *testing.T) {
	var output = `acl 2.2.52-2 amd64 Access control list utilities
adduser 3.113+nmu3 all add and remove users and groups
apt 1.0.9.8.2 amd64 commandline package manager
`
	pkgs := parseMockOutput(output)
	expected_pkgs := []pkg{
		pkg{
			name:         "acl",
			version:      "2.2.52-2",
			architecture: "amd64",
			summary:      "Access control list utilities",
		},
		pkg{
			name:         "adduser",
			version:      "3.113+nmu3",
			architecture: "all",
			summary:      "add and remove users and groups",
		},
		pkg{
			name:         "apt",
			version:      "1.0.9.8.2",
			architecture: "amd64",
			summary:      "commandline package manager",
		},
	}
	k := 0
	for p := range pkgs {
		assert.Equal(t, expected_pkgs[k], p)
		k++
	}
	assert.Equal(t, k, 3)
}

func TestParseOutputEmpty(t *testing.T) {
	pkgs := parseMockOutput("")
	_, ok := <-pkgs
	assert.False(t, ok)

	pkgs = parseMockOutput(" ")
	_, ok = <-pkgs
	assert.False(t, ok)
}
