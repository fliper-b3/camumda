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
	startedProc, err :=getStartedProc("http://192.168.64.2:31700/engine-rest/history/process-instance/count")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, fmt.Sprintf("camunda_processes_started_per_instance{namespace=\"defualt\"}  %d", startedProc/pods))
}

func getStartedProc(url string) (int, error) {
	client := http.Client{}
	rj := responseJSON{}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Print("ehlo drow")
		//return 0, err
	}
	json.NewDecoder(resp.Body).Decode(&rj)
	defer resp.Body.Close()
	return rj.Count, nil
}

func getPods(namespace string, labelKey, labelValue string) (int, error) {

	cln, err := newClientSet(true)

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
	http.HandleFunc("/metric", getMetric)
	http.ListenAndServe(":80", nil)
}