# Autoscale

This additional container in the camundaBSP and monitoring system which use custom metrics,
horizontal pod autoscaler and hpa polices to control workload

## Build

```bash
cd camunda-autoscaler/hpa/src
docker build -t avguston/metrics:stable
docker push avguston/metrics:stable
```


## Install

```bash
cd camunda-autoscaler/hpa
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prmth stable/prometheus-operator
helm install prometheus-adapter prometheus-community/prometheus-adapter
kubectl apply -f yaml/
```

## Clean

```bash
helm delete prmth
helm delete prometheus-adapter
kubectl delete deployment camunda-deployment
kubectl delete deployment camunda-process-starter
```

### Pros'n'Cons

+ controls by k8s api-server
+ k8s native
+ less custom code

- needs full production ready infrastraction
- needs to change camunda deployment