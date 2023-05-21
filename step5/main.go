package main

import (
	"devgo/step5/finance"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	f := finance.NewFinace("data")
	err := f.GetRankData()
	if err != nil {
		log.Error(err)
	}

}
