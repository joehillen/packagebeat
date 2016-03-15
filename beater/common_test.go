package beater

import (
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
