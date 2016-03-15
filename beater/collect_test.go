package beater

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLineFieldCount(t *testing.T) {
	assert.Nil(t, parseLine(""))
	assert.Nil(t, parseLine(" "))
	assert.Nil(t, parseLine("vim "))
	assert.Nil(t, parseLine("vim 2:7.4.052-1ubuntu3"))
	assert.NotNil(t, parseLine("vim 2:7.4.052-1ubuntu3 amd64"))
	assert.NotNil(t, parseLine("vim 2:7.4.052-1ubuntu3 amd64     "))
	assert.NotNil(t, parseLine("vim 2:7.4.052-1ubuntu3 amd64 Vi"))
	assert.NotNil(t, parseLine("vim 2:7.4.052-1ubuntu3 amd64 Vi "))
	assert.NotNil(t, parseLine("vim 2:7.4.052-1ubuntu3 amd64 Vi IMproved - enhanced vi editor"))
}

func TestParseLine(t *testing.T) {
	assert.Equal(t,
		pkg{
			name:         "vim",
			version:      "2:7.4.052-1ubuntu3",
			architecture: "amd64",
			summary:      "Vi ",
		},
		*parseLine("vim 2:7.4.052-1ubuntu3 amd64 Vi "),
	)
	assert.Equal(t,
		pkg{
			name:         "vim",
			version:      "2:7.4.052-1ubuntu3",
			architecture: "amd64",
			summary:      "Vi IMproved - enhanced vi editor",
		},
		*parseLine("vim 2:7.4.052-1ubuntu3 amd64 Vi IMproved - enhanced vi editor"),
	)
}

func TestParseLineNoDescription(t *testing.T) {
	assert.Equal(t,
		pkg{
			name:         "vim",
			version:      "2:7.4.052-1ubuntu3",
			architecture: "amd64",
			summary:      "",
		},
		*parseLine("vim 2:7.4.052-1ubuntu3 amd64"),
	)
}

func parseMockOutput(output string) chan pkg {
	return parseOutput(ioutil.NopCloser(strings.NewReader(output)))
}

func TestParseOutputEmpty(t *testing.T) {
	pkgs := parseMockOutput("")
	_, ok := <-pkgs
	assert.False(t, ok)

	pkgs = parseMockOutput(" ")
	_, ok = <-pkgs
	assert.False(t, ok)
}

func assertParseOutput(t *testing.T, output string, expected []pkg) {
	pkgs := parseMockOutput(output)
	k := 0
	for p := range pkgs {
		assert.Equal(t, expected[k], p)
		k++
	}
	assert.Equal(t, k, len(expected))
}

func TestDpkg(t *testing.T) {
	assertParseOutput(t,
		`acl 2.2.52-2 amd64 Access control list utilities
adduser 3.113+nmu3 all add and remove users and groups
apt 1.0.9.8.2 amd64 commandline package manager
`,
		[]pkg{
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
		},
	)
}

func TestRpm(t *testing.T) {
	assertParseOutput(t,
		`fedora-repos 23 noarch Fedora package repositories
basesystem 11 noarch The skeleton package which defines a simple Fedora system
libstdc++ 5.3.1 x86_64 GNU Standard C++ Library
`,
		[]pkg{
			pkg{
				name:         "fedora-repos",
				version:      "23",
				architecture: "noarch",
				summary:      "Fedora package repositories",
			},
			pkg{
				name:         "basesystem",
				version:      "11",
				architecture: "noarch",
				summary:      "The skeleton package which defines a simple Fedora system",
			},
			pkg{
				name:         "libstdc++",
				version:      "5.3.1",
				architecture: "x86_64",
				summary:      "GNU Standard C++ Library",
			},
		},
	)
}
