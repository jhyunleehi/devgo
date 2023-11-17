package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

type SoClient struct {
	Id        string
	Passwd    string
	SoIP      string
	SoPort    string
	Protocol  string
	transport *http.Transport
}

type ErrorMessage struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

const (
	HTTP_CLIENT_TIMEOUT        = 120
	HTTP_TRANSPORT_TLS_TIMEOUT = 5
	NET_DIALER_TIMEOUT         = 3
)

type UpdateVolumeOpts struct {
	StorageId        int64                  `json:"-"`
	VolumePoolId     string                 `json:"-"`
	VolumeName       string                 `json:"volumeName"`
	VolumeId         string                 `json:"-"`
	State            string                 `json:"state" example:"offline" enums:"online,offline"` // This Field Value is required, when mode=state
	Qos              VolumeQos              `json:"qos,omitempty"`
	TotalByte        int64                  `json:"totalByte" example:"1024"` // This Field Value is required, when mode=extend
	SpaceReservation string                 `json:"spaceReservation,omitempty" enums:"thin,thick"`
	Description      string                 `json:"description" example:"storage volume descriptions"`
	NFS              VolumeNFSOpts          `json:"nfs,omitempty"`
	CIFS             VolumeCIFSOpts         `json:"cifs,omitempty"`
	ISCSI            VolumeISCSIOpts        `json:"iscsi,omitempty"`
	Object           VolumeObjectOpts       `json:"object,omitempty"`
	RBD              VolumeRBDOpts          `json:"rbd,omitempty"`
	CustomParams     map[string]interface{} `json:"customParams"` //  example: { "transactionId":"TRC00001" }
}

type VolumeQos struct {
	MinThroughputIops int64 `json:"minThroughputIops" example:"0"`
	MaxThroughputIops int64 `json:"maxThroughputIops" example:"0"`
	MaxThroughputMbps int64 `json:"maxThroughputMbps" example:"0"`
}

type VolumeGroup struct {
	Name string `json:"name,omitempty" example:"volumeGroupName"`
}

type VolumeNFSOpts struct {
	Enable       string `json:"enable" example:"true"`
	ExportPolicy string `json:"exportPolicy" example:"192.168.57.4,192.168.57.0/24"`
}

type VolumeNFSSpec struct {
	Enable       string `json:"enable" example:"true"`
	ExportPolicy string `json:"exportPolicy" example:"192.168.57.4,192.168.57.0/24"`
	Mountpoint   string `json:"mountpoint" example:"192.168.57.4:/vol1"`
}

type VolumeCIFSShareUser struct {
	Id       string `json:"id" example:"userId"`
	Password string `json:"password" example:"password"`
	Acl      string `json:"acl" example:"FullControl"`
}

type VolumeCIFSOpts struct {
	Enable string                `json:"enable" example:"true"`
	User   []VolumeCIFSShareUser `json:"user,omitempty"`
}

type VolumeISCSIOpts struct {
	Enable    string                 `json:"enable" example:"true"`
	LunId     string                 `json:"lunId" example:"3600a09803830354e7924512f584d4a44"`
	PortalIp  string                 `json:"portalIp" example:"192.168.10.2,192.168.11.2"`
	Initiator []VolumeISCSIInitiator `json:"initiator"`
}

type VolumeISCSIInitiator struct {
	InitiatorKey  string `json:"initiatorKey" example:"192.168.56.3"`                                     // This Field Value is required only for ISCSI Storage.
	InitiatorName string `json:"initiatorName,omitempty" example:"iqn.1993-08.org.debian:01:3aece8b652e"` // This Field Value is required only for ISCSI Storage.
	InitiatorOS   string `json:"initiatorOS" example:"LINUX"`                                             //This Field Value is required only for ISCSI Storage. (LINUX, VMWARE, WINDOWS )
	InitiatorType string `json:"initiatorType" example:"IQN" enums:"IQN,NQN"`
}

type VolumeRBDOpts struct {
	Enable      string   `json:"enable" example:"true"`
	ObjectSize  int      `json:"objectsize,omitempty"`
	RbdFeatures []string `json:"features,omitempty"`
	Namespace   string   `json:"namespace,omitempty"`
}

type VolumeObjectOpts struct {
	Enable   string        `json:"enable" example:"true"`
	Endpoint string        `json:"endpoint"`
	Users    []StorageUser `json:"users"`
}

type StorageUser struct {
	UserId     string `json:"userId"`
	AccessKey  string `json:"accessKey"`
	Permission string `json:"permission" example:"disabled,readOnly,readWrite"`
}

type GetVolumeOpts struct {
	Mode                string
	StorageId           int64
	VolumeId            string
	VolumePoolId        string
	VolumePoolName      string
	VolumeName          string
	SelectVolumePoolIds []string
	CustomParams        map[string]interface{} `json:"customParams"` //  example: { "transactionId":"TRC00001" }

}

type VolumeSpec struct {
	StorageId            int64                  `json:"storageId" example:"1"`
	VolumePoolId         string                 `json:"volumePoolId" example:"SVM01:Aggr01"`
	VolumeName           string                 `json:"volumeName" example:"vol1"`
	VolumeId             string                 `json:"volumeId" example:"0e432327-b05b-11ea-986e-00a098638369"`
	VolumeGroup          VolumeGroup            `json:"volumeGroup,omitempty"`
	State                string                 `json:"state" enums:"online,offline"`
	Reserve              string                 `json:"reserve" example:"false"`
	Type                 string                 `json:"type" example:"rw"`
	Style                string                 `json:"style,omitempty" enums:"group"`
	Qos                  VolumeQos              `json:"qos,omitempty"`
	Encryption           string                 `json:"encryption" example:"false"`
	UsedByte             int64                  `json:"usedByte" example:"1024"`
	TotalByte            int64                  `json:"totalByte" example:"1024"`
	SpaceReservation     string                 `json:"spaceReservation,omitempty" enums:"thin,thick"`
	Honored              string                 `json:"honored,omitempty" example:"true/false"`
	Description          string                 `json:"description" example:"volume description"`
	NFS                  VolumeNFSSpec          `json:"nfs,omitempty"`
	CIFS                 VolumeCIFSOpts         `json:"cifs,omitempty"`
	ISCSI                VolumeISCSIOpts        `json:"iscsi,omitempty"`
	Object               VolumeObjectOpts       `json:"object,omitempty"`
	SubVolume            []SubVolumeSpec        `json:"subVolume,omitempty"`
	SnapshotSchedule     []SnapshotSchedule     `json:"snapshotSchedule"`
	SnapshotReservedRate int64                  `json:"snapshotReservedRate"`
	SnapshotPathVisible  string                 `json:"snapshotPathVisible" example:"false"`
	SnapshotUsedByte     int64                  `json:"snapshotUsedByte"`
	SnapshotAutoDelete   string                 `json:"snapshotAutoDelete" example:"true"`
	TaskId               string                 `json:"taskId,omitempty"`
	CustomParams         map[string]interface{} `json:"customParams"` //  example: { "transactionId":"TRC00001" }
}
type SubVolumeSpec struct {
	SubVolumeName string `json:"subVolumeName" example:"sub1"`
	Mountpoint    string `json:"mountpoint" example:"/vol70/qtree70-1"`
	TotalByte     int64  `json:"totalByte" example:"1024"`
	UsedByte      int64  `json:"usedByte" example:"1024"`
	ExportPolicy  string `json:"exportPolicy" example:"192.168.57.4,192.168.57.0/24"`
}

type SnapshotSchedule struct {
	SnapshotScheduleId   string `json:"snapshotScheduleId" example:"id"`
	SnapshotScheduleName string `json:"snapshotScheduleName" example:"name"`
}

func NewClient(ip, port, id, passwd, protocol string) (*SoClient, error) {
	log.Debugf("[%+v][%+v][%+v][%+v]", ip, port, id, passwd)

	c := &SoClient{
		Id:       id,
		Passwd:   passwd,
		SoIP:     ip,
		SoPort:   port,
		Protocol: protocol,
	}

	c.transport = &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: HTTP_TRANSPORT_TLS_TIMEOUT * time.Second,
		Dial:                (&net.Dialer{Timeout: NET_DIALER_TIMEOUT * time.Second}).Dial,
	}

	return c, nil
}

func (c *SoClient) doRequest(method, url string, in, out interface{}) (http.Header, error) {
	var inbody []byte
	var body *bytes.Buffer
	var req *http.Request
	if in != nil {
		inbody, _ = json.Marshal(in)
		body = bytes.NewBuffer(inbody)
		req, _ = http.NewRequest(method, url, body)
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Key", c.Id)
	req.Header.Add("X-Auth-Secret", c.Passwd)

	client := http.Client{
		Transport: c.transport,
		Timeout:   HTTP_CLIENT_TIMEOUT * time.Second,
	}

	resp, errReq := client.Do(req)
	if errReq != nil {
		if resp != nil {
			buf1, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Print(err.Error())
				return nil, err
			}
			fmt.Printf("[%+v]", string(buf1))
			log.Debugf("get bytes[] from response error: [%s] [%s] [%+v] [%+v]", method, url, errReq, string(buf1))
			return nil, errReq
		}
		return nil, errReq
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		buf, _ := ioutil.ReadAll(resp.Body)
		APIError := &ErrorMessage{}
		json.Unmarshal([]byte(buf), APIError)
		//if APIError.Code == 6 && sessionRetry == false {
		log.Errorf("[%+v][%+v]", resp.Status, resp.StatusCode)

	}
	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("http response status code: [%d]", resp.StatusCode)
		log.Errorf(msg)
		buf, _ := ioutil.ReadAll(resp.Body)
		log.Errorf(string(buf))
		return http.Header{}, errors.New(string(buf))
	}

	msg := fmt.Sprintf("http response status code: [%d]", resp.StatusCode)
	log.Debug(msg)

	buf, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 && method != "GET" {
		msg := fmt.Sprintf("Resource Updated [%v]", resp.StatusCode)
		log.Debug(msg)
		log.Debugf(string(buf))
	}
	if resp.StatusCode == 201 {
		msg := fmt.Sprintf("Resource Created [%v]", resp.StatusCode)
		log.Debug(msg)
		log.Debugf(string(buf))
	}
	if resp.StatusCode == 202 {
		msg := "Operation is still executing. Please check the task queue."
		log.Debug(msg)
		log.Debugf(string(buf))
	}
	if resp.StatusCode == 204 {
		msg := fmt.Sprintf("Resource deleted. [%v]", resp.StatusCode)
		log.Debug(msg)
		log.Debugf(string(buf))
	}

	if out != nil {
		err := json.Unmarshal([]byte(buf), out)
		if err != nil {
			log.Debug(string(buf))
			log.Error(err)
		}
	}
	return resp.Header, nil
}

func (c *SoClient) GetVolume(req *GetVolumeOpts, storageid string) (*[]VolumeSpec, error) {
	log.Debugf("[%+v]", req)
	vollist := &[]VolumeSpec{}
	url := fmt.Sprintf("%s://%s:%s/api/v1/storages/%s/volume-pools/%s/volumes/?volumeName=%s", c.Protocol, c.SoIP, c.SoPort, storageid, req.VolumePoolId, req.VolumeName)
	log.Debugf("[%s]", url)
	_, errReq := c.doRequest("GET", url, req, vollist)
	if errReq != nil {
		log.Debugf("GetVolume error [%+v]", errReq)
		return vollist, errReq
	}
	return vollist, nil
}

func Get_VolumeId_by_Name(poolid, volname, sid string) (string, error) {
	var volumeId string
	opt := GetVolumeOpts{}
	opt.VolumePoolId = poolid
	opt.VolumeName = volname
	rtn, err := c.GetVolume(&opt, sid)
	if err != nil {
		log.Errorf("[%+v][%+v]", err, rtn)
		return "", err
	}
	if len(*rtn) > 0 {
		volumeId = (*rtn)[0].VolumeId
	}
	return volumeId, nil
}

var c *SoClient
var err error
var IP, PORT, ID, AUTH, sid string

func init() {
	IP = "70.60.31.51"
	PORT = "10443"
	ID = "sds"
	AUTH = "1234"
	sid = "231"
	c, err = NewClient(IP, PORT, ID, AUTH, "https")
	if c == nil || err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	log.SetLevel(log.InfoLevel)

	//log.SetFormatter(&log.TextFormatter{DisableColors: false})
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		TimestampFormat: time.RFC3339,
		NoColors:        true,
	})
}

func main() {
	modeADD := "?mode=nfsExportPolicyADD"
	modeMOD := "?mode=nfsExportPolicyMOD"
	//VOLUME_POOL_NAME := "POOL_NFS"
	VOLUME_POOL_ID_NFS := "POOL_NFS:N1_aggr1"
	VOLUME_NAME_NFS := "VOL_NFS"

	VOLUME_ID, err := Get_VolumeId_by_Name(VOLUME_POOL_ID_NFS, VOLUME_NAME_NFS, sid)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	for {
		vol := &VolumeSpec{}
		opt := UpdateVolumeOpts{}
		opt.VolumePoolId = VOLUME_POOL_ID_NFS
		opt.VolumeId = VOLUME_ID
		opt.NFS.ExportPolicy = ""
		url := fmt.Sprintf("%s://%s:%s/api/v1/storages/%s/volume-pools/%s/volumes/%s%s", c.Protocol, c.SoIP, c.SoPort, sid, opt.VolumePoolId, opt.VolumeId, modeMOD)
		log.Debugf("[%s]", url)
		_, errReq := c.doRequest("PUT", url, opt, vol)
		if errReq != nil {
			log.Errorf("[%+v]", err)
			log.Panic()
			return
		}

		for k := 1; k <= 2; k++ {
			for i := 1; i <= 100; i++ {
				opt.NFS.ExportPolicy = "1.1." + strconv.Itoa(k) + "." + strconv.Itoa(i)
				url := fmt.Sprintf("%s://%s:%s/api/v1/storages/%s/volume-pools/%s/volumes/%s%s", c.Protocol, c.SoIP, c.SoPort, sid, opt.VolumePoolId, opt.VolumeId, modeADD)
				log.Debugf("[%s]", url)
				_, errReq := c.doRequest("PUT", url, opt, vol)
				if errReq != nil {
					log.Errorf("[%+v]", err)
					log.Panic()
					return
				}

				o := GetVolumeOpts{}
				o.VolumePoolId = VOLUME_POOL_ID_NFS
				o.VolumeName = VOLUME_NAME_NFS
				getrtn, err := c.GetVolume(&o, sid)
				if err != nil {
					log.Errorf("[%+v]", err)
					log.Panic()
					return
				}

				networklist := strings.Split((*getrtn)[0].NFS.ExportPolicy, ",")
				log.Printf("[%d]  [%d]", (k-1)*100+i, len(networklist))
				if len(networklist) != (k-1)*100+i {
					log.Errorf("[%+v]", networklist)
					log.Errorf("=====>>>>>[%+v]", networklist)
					log.Panic()
					return
				}
			}
		}
	}

}
