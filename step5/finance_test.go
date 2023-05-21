package main

import (
	"devgo/step5/finance"
	"testing"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
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

type TSuite struct {
	suite.Suite
	finance *finance.Finance
}

func Test_Suite(t *testing.T) {
	suite.Run(t, new(TSuite))
	log.Debugf("start test...")
}

func (s *TSuite) Test_Setup() {
	s.finance = finance.NewFinace("trend")
	s.finance.GetRankData()
}
