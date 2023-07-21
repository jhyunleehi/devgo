package main

import (
	"devgo/trend"
	"testing"

	"github.com/stretchr/testify/suite"	
	//log "github.com/sirupsen/logrus"
)

type TSuite struct {
	suite.Suite
	mytrend *trend.Trend
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TSuite))
}

func (s *TSuite) TestSetup() {	
	s.mytrend = trend.NewTrend("trend")
	s.mytrend.GetInit()
}

