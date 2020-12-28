package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func updateDeployment(replicasDelta int32, cln *kubernetes.Clientset, namespace, deploymentName string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deploy, err := cln.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "main",
				"function": "updateDeployment",
				"error":    err,
			}).Warningf("Cannot read a deployment %s", deploymentName)
			return err
		}
		*deploy.Spec.Replicas = *deploy.Spec.Replicas + replicasDelta
		_, err = cln.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "main",
				"function": "updateDeployment",
				"error":    err,
			}).Warningf("Cannot update the deployment %s", deploymentName)
			return err
		}
		return nil
	})

}

func autoScaleRules(pods, procStarted int, maxPod int) (int32, error) {
	if pods == 0 {
		return 1, errors.New("zero pod")
	}
	if procStarted/pods >= 50 && pods < maxPod {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "autoScaleRules",
		}).Debugf("Increase procStarted: %d pods: %d", procStarted, pods)
		return 1, nil
	}
	if procStarted/pods <= 20 && pods > 1 {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "autoScaleRules",
		}).Debugf("Decrease procStarted: %d pods: %d", procStarted, pods)
		return -1, nil
	}
	return 0, nil
}

func newClientSet(runOutsideCluster bool) (*kubernetes.Clientset, error) {

	config, err := rest.InClusterConfig()
	if err != nil && runOutsideCluster == false {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "newClientSet",
			"error":    err,
		}).Fatal("Cannot read a kube config")
		return nil, err
	}
	if runOutsideCluster == true {
		kubeConfigLocation := ""
		homeDir := os.Getenv("HOME")
		kubeConfigLocation = filepath.Join(homeDir, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigLocation)
	}
	return kubernetes.NewForConfig(config)
}

func getStartedProc(url string) (int, error) {
	type responseJSON struct {
		Count int
	}
	rj := responseJSON{}
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode == 200 {
		json.NewDecoder(resp.Body).Decode(&rj)
		return rj.Count, nil
	}
	defer resp.Body.Close()
	return 0, errors.New("Respons code error")
}

func getPods(namespace string, labelKey, labelValue string, cln *kubernetes.Clientset, maxPod int) (int, error) {
	var readyPodsCount int
	listOptions := metav1.ListOptions{LabelSelector: fmt.Sprintf("%s=%s", labelKey, labelValue)}
	podsList, err := cln.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	if err != nil {
		return 0, err
	}
	totalPods := len(podsList.Items)

	if totalPods >= maxPod {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "getPods",
		}).Warning("Maximum pods exist")
		return totalPods, nil
	}

	for i := 0; i < totalPods; i++ {
		if *podsList.Items[i].Status.ContainerStatuses[0].Started {
			readyPodsCount++
		}
	}

	return readyPodsCount, nil
}

func work(cln *kubernetes.Clientset, url, namespace, labelKey, labelValue, deploymentName string, maxPod int) {
	now := time.Now().UTC()
	tenSecond := now.Add(-10 * time.Second)
	pods, err := getPods(namespace, labelKey, labelValue, cln, maxPod)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "work",
			"error":    err,
		}).Error("Cannot get pods")
		return
	}
	procStrated, err := getStartedProc(fmt.Sprintf(
		"%s?startedAfter=%s", url,
		tenSecond.Format("2006-01-02T15:04:05.56-0000"),
	))
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "work",
			"error":    err,
		}).Errorf(
			"Cannot read url %s", url,
		)
		return
	}
	scale, err := autoScaleRules(pods, procStrated, maxPod)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "work",
		}).Error(err)
	}
	if scale != 0 && err == nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "work",
		}).Infof(
			"Updated the deployment %s, procStarted: %d, new pods count: %d",
			deploymentName, procStrated, scale+int32(pods),
		)
		err = updateDeployment(scale, cln, namespace, deploymentName)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "main",
				"function": "work",
				"error":    err,
			}).Errorf("Cannot updated the deployment %s", deploymentName)
			return
		}
	}
}

func main() {
	runOutsideCluster := false
	namespace := "default"
	deploymentName := "camunda-deployment"
	labelKey := "app"
	labelValue := "camunda"
	countUrl := "http://camunda-service:8080/engine-rest/history/process-instance/count"
	log.SetLevel(log.InfoLevel)
	cln, err := newClientSet(runOutsideCluster)
	if runOutsideCluster {
		log.SetLevel(log.DebugLevel)
		countUrl = "http://127.0.0.1:57093/engine-rest/history/process-instance/count"
	}
	if err != nil {
		panic(err)
	}
	log.Infof(
		"Starting...\n\tcamunda url: %s\n\tnamespace:%s\n\tdeployment name: %s",
		countUrl, namespace, deploymentName,
	)
	for {
		work(
			cln,
			countUrl,
			namespace,
			labelKey,
			labelValue,
			deploymentName,
			4,
		)
		sec := 10 * time.Second
		log.Debugf("sleeping %d sec", sec)
		time.Sleep(sec)
	}

}
