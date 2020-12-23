package main

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type responseJSON struct {
	Count int
}

func newClientSet(runOutsideCluster bool) (*kubernetes.Clientset, error) {

	config, err := rest.InClusterConfig()

	if runOutsideCluster == true {
		kubeConfigLocation := ""
		homeDir := os.Getenv("HOME")
		kubeConfigLocation = filepath.Join(homeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigLocation)
	}

	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func getMetric (w http.ResponseWriter, r *http.Request) {
	pods, err:=getPods("default", "app", "camunda")
	if err != nil {
		panic(err)
	}
	now := time.Now().UTC()
	tenSecond := now.Add(-10*time.Second)
	startedProc, err :=getStartedProc(fmt.Sprintf(
		"http://localhost:8080/engine-rest/history/process-instance/count?startedAfter=%s",
		tenSecond.Format("2006-01-02T15:04:05-0700"),
	))
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, fmt.Sprintf("camunda_processes_started_per_instance %d", startedProc/pods))
	fmt.Println(fmt.Sprintf("camunda_processes_started_per_instance %d", startedProc/pods))
}

func getStatus (w http.ResponseWriter, r *http.Request) {
	fmt.Println("hello word")
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

func getPods(namespace string, labelKey, labelValue string) (int, error) {

	cln, err := newClientSet(false)

	if err != nil {
		panic(err.Error())
	}
	listOptions := metav1.ListOptions{LabelSelector:fmt.Sprintf("%s=%s", labelKey, labelValue)}
	podsList, err := cln.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if  err != nil {
		return 0, err
	}
	return len(podsList.Items), nil
}

func main() {
	http.HandleFunc("/metrics", getMetric)
	http.HandleFunc("/status", getStatus)
	http.ListenAndServe(":8088", nil)
}