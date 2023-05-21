package main

import (	
	"sync"
	"time"
	
	resty "github.com/go-resty/resty/v2"
	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
		NoColors:        true,
	})
}

type Resty struct {
	name string
	client  
	mutex sync.Mutex
}

func NewResy(name string) *Resty {
	
	return &Resty{
		name: name,
	}
	Resty.client = resty.New()
	return resty
}
