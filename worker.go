package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type WorkerSpec struct {
	host     string
	login    string
	password string
	pause    uint
}

type Worker struct {
	sync.Mutex
	spec   WorkerSpec
	target *Target
}

func NewWorker(spec WorkerSpec) *Worker {
	log.Println("New worker (login:", spec.login, ", host:", spec.host, ")")

	// create Worker
	w := Worker{
		spec:   spec,
		target: NewTarget(spec.host),
	}

	// start main loop
	go w.getPIRStatus(1000 * time.Millisecond)

	return &w
}

func (w *Worker) GetWorkerTarget(host string) *Target {
	w.Lock()
	defer w.Unlock()
	t := w.target
	if t.host == "" {
		t = NewTarget(host)
		w.target = t
	}
	return t
}

func (w *Worker) getPIRStatus(sleepTime time.Duration) {
	time.Sleep(sleepTime)

	go w.getPIRStatus(time.Duration(w.spec.pause) * time.Second)

	uri := fmt.Sprintf("http://%s:%s@%s/ISAPI/WLAlarm/PIR", w.spec.login, w.spec.password, w.spec.host)
	client := &http.Client{}
	client.Timeout = time.Second * 2
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("Content-Type", "text/xml")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("resp, err := client.Do(req): ", err)
		return
	}
	defer resp.Body.Close()
	xmlFile, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("xmlFile, err := ioutil.ReadAll(resp.Body): ", err)
		return
	}
	pirAlarmXML := PIRAlarmXML{}
	err = xml.Unmarshal(xmlFile, &pirAlarmXML)
	if err != nil {
		log.Println("err = xml.Unmarshal(xmlFile, &pirAlarmXML)", err)
		return
	}
	w.addResults(pirAlarmXML)
}

func (w *Worker) addResults(pirAlarmXML PIRAlarmXML) {
	// Parse results
	pirAlarmStatus, err := ParsePIRAlarmStatus(pirAlarmXML)
	if err != nil {
		log.Println("Error parsing xml: ", err)
		return
	}
	log.Println(pirAlarmStatus)
	w.target.AddPIRStatus(pirAlarmStatus)
}
