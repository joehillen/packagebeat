package beater

import (
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
)

type Packagebeat struct {
	period   time.Duration
	PbConfig ConfigSettings
	events   publisher.Client
	done     chan struct{}
}

func New() *Packagebeat {
	return &Packagebeat{}
}

func (pb *Packagebeat) Config(b *beat.Beat) error {
	err := cfgfile.Read(&pb.PbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	if pb.PbConfig.Input.Period == nil {
		pb.period = 60 * time.Second
	} else {
		pb.period = time.Duration(*pb.PbConfig.Input.Period) * time.Second
	}

	logp.Debug("packagebeat", "Init packagebeat")
	logp.Debug("packagebeat", "Period %v\n", pb.period)

	return nil
}

func (pb *Packagebeat) Setup(b *beat.Beat) error {
	pb.events = b.Events
	pb.done = make(chan struct{})
	return nil
}

func (pb *Packagebeat) Run(b *beat.Beat) error {
	for {
		pb.CollectDpkg()
		pb.CollectRPM()
		time.Sleep(pb.period)
	}
}

func (pb *Packagebeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (pb *Packagebeat) Stop() {
	close(pb.done)
}
