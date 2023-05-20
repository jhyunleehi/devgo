package finance

import (	
	"sync"
	"time"

	"github.com/anaskhan96/soup"
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

type Finance struct {
	name  string
	mutex sync.Mutex
}

func NewFinace(name string) *Finance {
	return &Finance{
		name: name,
	}
}

func (f *Finance) GetRankData() error {
	resp, err := soup.Get("https://finance.naver.com/sise/lastsearch2.naver")
	if err != nil {
		log.Error(err)
		return err
	}
	//fmt.Printf("%s", resp)
	doc := soup.HTMLParse(resp)
	div := doc.FindAll("tbody")//("div", "class", "box_type_1")
	for _, d := range div {
		links := d.FindAll("tr")
		for _, link := range links {
			rankitem := link.Text()
			//rankitem := (link.Attrs()["title"])			
			log.Debugf("[%s]",rankitem)
			//fmt.Println(link.Text(), "| Link :", link.Attrs()["href"])
		}
	}
	return nil
}
