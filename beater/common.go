package beater

import (
	"regexp"
	"strings"

	"github.com/elastic/beats/libbeat/logp"
)

type pkg struct {
	name         string
	version      string
	architecture string
	summary      string
}

func parseLine(line string) *pkg {
	words := regexp.MustCompile("\\s+").Split(line, 4)
	cnt := len(words)
	if cnt < 3 {
		logp.Err("Not enough fields (%d) in line: %s", len(words), line)
		return nil
	}

	var summary = ""
	if cnt > 3 {
		summary = strings.Join(words[3:], " ")
	}

	return &pkg{
		name:         words[0],
		version:      words[1],
		architecture: words[2],
		summary:      summary,
	}
}
