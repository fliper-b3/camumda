# Autoscale

This additional container in the camundaBSP and monitoring system which use custom metrics,
horizontal pod autoscaler and hpa polices to control workload

## Build

```bash
cd camunda-autoscaler/hpa
docker build -t <reponame>/metrics:stable .
docker push <reponame>/metrics:stable
```

## Install

```bash
cd camunda-sre-interview-master/
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prmth stable/prometheus-operator
helm install prometheus-adapter prometheus-community/prometheus-adapter \
  --set prometheus.url=http://prmth-prometheus-operator-prometheus
kubectl apply -f k8s_resources/postgres.yml
kubectl apply -f camunda-autoscaler/hpa/yaml/
kubectl apply -f k8s_resources/processStarter.yml
```

## Manage a payload

To increase/decrease payload change in a file `camunda-sre-interview-master/k8s_resources/processStarter.yml`
value for env variable `N_PROCESS_STARTED` or `QUIET_TIME_S` or both

## Clean

```bash
helm delete prmth
helm delete prometheus-adapter
kubectl delete -f k8s_resources/postgres.yml
kubectl delete -f camunda-autoscaler/hpa/yaml/
kubectl delete -f k8s_resources/processStarter.yml
```

## Pros'n'Cons

(+) controls by k8s api-server

(+) k8s native

(+) less custom code

(+) production ready more or less

(-) needs full production ready infrastructure

(-) needs to change camunda deployment or container
