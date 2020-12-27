package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func getMetric (w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	tenSecond := now.Add(-10*time.Second)
	url := fmt.Sprintf(
		"http://localhost:8080/engine-rest/history/process-instance/count?startedAfter=%s",
		tenSecond.Format("2006-01-02T15:04:05.56-0000"),
	)
	startedProc, err :=getStartedProc(url)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "getMetric",
			"error":    err,
		}).Errorf("Cannot read the url %s", url)
	}
	fmt.Fprintf(w, fmt.Sprintf("camunda_processes_started_per_instance %d", startedProc))
	log.WithFields(log.Fields{
		"package":  "main",
		"function": "getMetric",
	}).Infof(fmt.Sprintf("{\"request\": \"%s\", \"respons\": \"camunda_processes_started_per_instance %d\"}", url, startedProc))
}

func getStartedProc(url string) (int, error) {
	client := http.Client{}
	type responseJSON struct {
		Count int
	}
	rj := responseJSON{}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	json.NewDecoder(resp.Body).Decode(&rj)
	defer resp.Body.Close()
	return rj.Count, nil
}

func getStatus (w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"package":  "main",
		"function": "getStatus",
	}).Debug("status checked")
	fmt.Fprintf(w, "I'm alive")
}

func main() {
	http.HandleFunc("/metrics", getMetric)
	http.HandleFunc("/status", getStatus)
	http.ListenAndServe(":8088", nil)
}