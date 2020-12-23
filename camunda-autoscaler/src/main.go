package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type responseJSON struct {
	Count int
}

func getMetric (w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	tenSecond := now.Add(-10*time.Second)
	startedProc, err :=getStartedProc(fmt.Sprintf(
		"http://localhost:8080/engine-rest/history/process-instance/count?startedAfter=%s",
		tenSecond.Format("2006-01-02T15:04:05-0700"),
	))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, fmt.Sprintf("camunda_processes_started_per_instance %d", startedProc))
	fmt.Println(fmt.Sprintf("camunda_processes_started_per_instance %d", startedProc))
}

func getStatus (w http.ResponseWriter, r *http.Request) {
	fmt.Println("status logs")
	fmt.Fprintf(w, "I'm alive")
}

func getStartedProc(url string) (int, error) {
	client := http.Client{}
	rj := responseJSON{}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	json.NewDecoder(resp.Body).Decode(&rj)
	defer resp.Body.Close()
	return rj.Count, nil
}

func main() {
	http.HandleFunc("/metrics", getMetric)
	http.HandleFunc("/status", getStatus)
	http.ListenAndServe(":8088", nil)
}