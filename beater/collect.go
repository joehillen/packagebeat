package beater

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
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

func parseOutput(output io.ReadCloser) chan pkg {
	pkgs := make(chan pkg)
	scanner := bufio.NewScanner(output)
	go func() {
		for scanner.Scan() {
			pkg := parseLine(scanner.Text())
			if pkg != nil {
				pkgs <- *pkg
			}
		}
		close(pkgs)
	}()
	return pkgs
}

func (pb *Packagebeat) collectPackages(manager string, cmd string, args ...string) error {
	now := common.Time(time.Now())
	x := exec.Command(cmd, args...)
	output, err := x.StdoutPipe()
	if err != nil {
		logp.Err("%v", err)
		return err
	}
	if err := x.Start(); err != nil {
		logp.Err("%v", err)
		return err
	}
	for pkg := range parseOutput(output) {
		pb.events.PublishEvent(common.MapStr{
			"@timestamp":   now,
			"type":         "package",
			"manager":      manager,
			"name":         pkg.name,
			"version":      pkg.version,
			"architecture": pkg.architecture,
			"summary":      pkg.summary,
		})
	}
	return nil
}

const DPKG_PATH = "/usr/bin/dpkg-query"

func (pb *Packagebeat) CollectDpkg() error {
	if _, err := os.Stat(DPKG_PATH); err == nil {
		logp.Info("packagebeat", "Collecting from dpkg")
		return pb.collectPackages("dpkg",
			DPKG_PATH, "--show", "--showformat",
			"${Package} ${Version} ${Architecture} ${binary:Summary}\n")
	}
	return nil
}

const RPM_PATH = "/usr/bin/rpm"

func (pb *Packagebeat) CollectRPM() error {
	if _, err := os.Stat(RPM_PATH); err == nil {
		logp.Info("packagebeat", "Collecting from rpm")
		return pb.collectPackages("rpm",
			RPM_PATH, "-qa", "--queryformat",
			"%{NAME} %{VERSION} %{ARCH} %{SUMMARY}\n")
	}
	return nil
}
