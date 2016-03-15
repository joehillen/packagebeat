package beater

import (
	"bufio"
	"io"
	"os/exec"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

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

func (pb *Packagebeat) CollectDpkg() error {
	logp.Info("packagebeat", "Collecting packages from dpkg")
	now := common.Time(time.Now())
	dpkg := exec.Command(
		"/usr/bin/dpkg-query", "--show", "--showformat",
		"${Package} ${Version} ${Architecture} ${binary:Summary}\n")
	dpkgOutput, err := dpkg.StdoutPipe()
	if err != nil {
		logp.Err("%v", err)
		return err
	}
	if err := dpkg.Start(); err != nil {
		logp.Err("%v", err)
		return err
	}
	for pkg := range parseOutput(dpkgOutput) {
		pb.events.PublishEvent(common.MapStr{
			"@timestamp":   now,
			"type":         "package",
			"manager":      "dpkg",
			"name":         pkg.name,
			"version":      pkg.version,
			"architecture": pkg.architecture,
			"summary":      pkg.summary,
		})
	}
	return nil
}
