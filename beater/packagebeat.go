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

	dpkg bool
	rpm  bool
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

	if pb.PbConfig.Input.Dpkg != nil {
		pb.dpkg = *pb.PbConfig.Input.Dpkg
	} else {
		pb.dpkg = true
	}

	if pb.PbConfig.Input.Rpm != nil {
		pb.rpm = *pb.PbConfig.Input.Rpm
	} else {
		pb.rpm = true
	}

	logp.Debug("packagebeat", "Init packagebeat")
	logp.Debug("packagebeat", "Period %v\n", pb.period)
	logp.Debug("packagebeat", "rpm %t\n", pb.rpm)
	logp.Debug("packagebeat", "dpkg %t\n", pb.dpkg)

	return nil
}

func (pb *Packagebeat) Setup(b *beat.Beat) error {
	pb.events = b.Publisher.Connect()
	pb.done = make(chan struct{})
	return nil
}

func (pb *Packagebeat) Run(b *beat.Beat) error {
	for {
		if pb.dpkg {
			pb.CollectDpkg()
		}
		if pb.rpm {
			pb.CollectRPM()
		}
		time.Sleep(pb.period)
	}
}

func (pb *Packagebeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (pb *Packagebeat) Stop() {
	close(pb.done)
}
