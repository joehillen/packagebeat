package beater

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseName(t *testing.T) {
	var derp = "derp"
	assert.Equal(t, parseName(derp), derp)
	assert.Equal(t, parseName(derp+":amd64"), derp)
	assert.Equal(t, parseName(derp+":"), derp)
	assert.Equal(t, parseName(derp+"::"), derp)
	assert.Equal(t, parseName(derp+"::foobar"), derp)
}

func TestParseLineFieldCount(t *testing.T) {
	assert.Nil(t, parseLine(""))
	assert.Nil(t, parseLine("ii"))
	assert.Nil(t, parseLine("ii  vim "))
	assert.Nil(t, parseLine("ii  vim 2:7.4.052-1ubuntu3"))
	assert.NotNil(t, parseLine("ii vim 2:7.4.052-1ubuntu3 amd64"))
	assert.NotNil(t, parseLine("ii vim 2:7.4.052-1ubuntu3 amd64     "))
	assert.NotNil(t, parseLine("ii vim 2:7.4.052-1ubuntu3 amd64 Vi IMproved - enhanced vi editor"))
}

func TestParseLineNotInstalled(t *testing.T) {
	assert.Nil(t, parseLine("un vim 2:7.4.052-1ubuntu3 amd64 Vi IMproved - enhanced vi editor"))
}

func TestParseLine(t *testing.T) {
	assert.Equal(t,
		*parseLine("ii vim 2:7.4.052-1ubuntu3 amd64 Vi IMproved - enhanced vi editor"),
		dpkgPackage{
			name:         "vim",
			version:      "2:7.4.052-1ubuntu3",
			architecture: "amd64",
			description:  "Vi IMproved - enhanced vi editor",
		},
	)

}

func TestParseLineNoDescription(t *testing.T) {
	assert.Equal(t,
		*parseLine("ii vim 2:7.4.052-1ubuntu3 amd64"),
		dpkgPackage{
			name:         "vim",
			version:      "2:7.4.052-1ubuntu3",
			architecture: "amd64",
			description:  "",
		},
	)
}

func mockOutput(output string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(output))
}

func TestParseOutput(t *testing.T) {
	var output = `Desired=Unknown/Install/Remove/Purge/Hold
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)
||/ Name                     Version                  Architecture Description
+++-========================-========================-============-======================================================================
ii  acl            2.2.52-1     amd64        Access control list utilities`
	pkgs := parseOutput(mockOutput(output))
	pkg, ok := <-pkgs
	assert.True(t, ok)
	assert.Equal(t, pkg.name, "acl")
	assert.Equal(t, pkg.version, "2.2.52-1")
	assert.Equal(t, pkg.architecture, "amd64")
	assert.Equal(t, pkg.description, "Access control list utilities")
	_, ok = <-pkgs
	assert.False(t, ok)
}

func TestParseOutputEmpty(t *testing.T) {
	pkgs := parseOutput(mockOutput(""))
	_, ok := <-pkgs
	assert.False(t, ok)
}
