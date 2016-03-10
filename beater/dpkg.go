package beater

import (
	"bufio"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

const DPKG_PATH = "/usr/bin/dpkg-query"

type dpkgPackage struct {
	name         string
	version      string
	architecture string
	description  string
}

func parseName(name string) string {
	return strings.Split(name, ":")[0]
}

func parseLine(line string) *dpkgPackage {

	words := regexp.MustCompile("\\s+").Split(line, 5)
	if len(words) < 4 {
		logp.Err("Not enough fields (%d) in dpkg line: %s", len(words), line)
		return nil
	}

	if words[0] != "ii" {
		return nil
	}

	var description = ""
	if len(words) > 4 {
		description = strings.Join(words[4:], " ")
	}

	return &dpkgPackage{
		name:         parseName(words[1]),
		version:      words[2],
		architecture: words[3],
		description:  description,
	}

}

func parseOutput(output io.ReadCloser) chan dpkgPackage {
	pkgs := make(chan dpkgPackage)
	scanner := bufio.NewScanner(output)
	var cnt = 0
	go func() {
		for scanner.Scan() {
			// skip first 4 lines of the header lines
			if cnt < 4 {
				cnt += 1
				continue
			}

			pkg := parseLine(scanner.Text())
			if pkg != nil {
				pkgs <- *pkg
			}
		}
		close(pkgs)
	}()
	return pkgs
}

func (pb *Packagebeat) CollectDpkg() error {
	logp.Debug("packagebeat", "Collection packages")

	dpkg := exec.Command(DPKG_PATH, "--list")
	dpkgOutput, err := dpkg.StdoutPipe()
	if err != nil {
		logp.Err("%v", err)
		return err
	}
	if err := dpkg.Start(); err != nil {
		logp.Err("%v", err)
		return err
	}

	now := common.Time(time.Now())

	pkgs := parseOutput(dpkgOutput)

	for pkg := range pkgs {
		pb.events.PublishEvent(common.MapStr{
			"@timestamp":   now,
			"type":         "package",
			"manager":      "dpkg",
			"name":         pkg.name,
			"version":      pkg.version,
			"architecture": pkg.architecture,
			"description":  pkg.description,
		})
	}
	return nil
}
