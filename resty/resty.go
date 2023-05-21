package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	resty "github.com/go-resty/resty/v2"
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
	name    string
	client  *resty.Client
	request *resty.Request
	mutex   sync.Mutex
}

func NewResty(name string) *Resty {
	r := Resty{
		name: name,
	}
	r.client = resty.New()
	//r.client.SetDebug(true)
	r.client.SetDebug(false)
	r.client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	r.client.SetTimeout(1 * time.Minute)
	r.request = r.client.R()
	return &r
}

func (r *Resty) RestApi(restitem interface{}) error {
	//log.Debugf("[%+v]", restitem)
	item := restitem.(map[string]interface{})

	name, ok := item["name"]
	if !ok {
		msg := fmt.Sprintf("Cannot find name filed [%s]", name)
		log.Debug(msg)
		return nil
	}
	log.Printf("Test:[%s]", name)

	req, ok := item["request"]
	if !ok {
		//return nil
		msg := fmt.Sprintf("Cannot find request filed [%s]", req)
		log.Debug(msg)
		return nil
	}
	request := req.(map[string]interface{})
	headers, ok := request["header"]
	if ok {
		for _, head := range headers.([]interface{}) {
			h := head.(map[string]interface{})
			r.request = r.request.SetHeader(h["key"].(string), h["value"].(string))
		}
	}
	url, ok := request["url"]
	if !ok {
		msg := fmt.Sprintf("Cannot find url field [%s]", url)
		log.Error(msg)
		return errors.New(msg)
	}
	urlmap := url.(map[string]interface{})
	raw := urlmap["raw"].(string)
	log.Debugf("%s", raw)

	body, ok := request["body"]
	if ok {
		bodymap := body.(map[string]interface{})
		bodystring := bodymap["raw"].(string)
		r.request.SetBody(bodystring)
	}

	method, ok := request["method"].(string)
	if !ok {
		msg := fmt.Sprintf("Cannot find method filed [%s]", method)
		log.Error(msg)
		return errors.New(msg)
	}
	switch method {
	case "GET":
		res, err := r.request.Get(raw)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(res.Body()) > 0 {
			log.Debugf("%s", string(res.Body()))
		}
	case "POST":
		res, err := r.request.Post(raw)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(res.Body()) > 0 {
			log.Debugf("%s", string(res.Body()))
		}
	case "PUT":
		res, err := r.request.Put(raw)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(res.Body()) > 0 {
			log.Debugf("%s", string(res.Body()))
		}
	case "PATCH":
		res, err := r.request.Patch(raw)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(res.Body()) > 0 {
			log.Debugf("%s", string(res.Body()))
		}
	case "DELETE":
		res, err := r.request.Delete(raw)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(res.Body()) > 0 {
			log.Debugf("%s", string(res.Body()))
		}
	default:

	}
	return nil
}
