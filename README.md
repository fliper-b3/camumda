# SRE Coding Challenge

Main description are stored in the `CodingChallengeTask.txt`

Details is finded in `camunda-autoscale/default/README.md`, `camunda-autoscaler/hpa/README.md` and `camunda-autoscaler/Questions.md`

### Userfull commands and articls

* To find your metric in k8s api `kubectl get --raw="/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/camunda_processes_started_per_instance" | jq`

* Get your metric from API kubernetes `kubectl proxy --port=8080` `http http://localhost:8080/api/v1/namespaces/default/services/camunda-metrics:web/proxy/metrics`

* Networktools `kubectl run tools --rm -it --image=travelping/nettools -- bash`

### Links 

[prometheus-adapter github](https://github.com/helm/charts/tree/master/stable/prometheus-adapter)

[hpa articl was very userfull](https://rtfm.co.ua/kubernetes-horizontalpodautoscaler-obzor-i-primery/#Applicationbased_metrics_scaling)

[minikube docs](https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/)
