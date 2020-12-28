# Autoscale

This is autoscaler based on client-go library and works by helm deployment

## Build

```bash
cd camunda-autoscaler/default
docker build -t <reponame>/camunda:stable .
docker push <reponame>/camunda:stable
```


## Install

```bash
cd camunda-sre-interview-master/
kubectl apply -f k8s_resources
cd camunda-autoscaler/default
helm install camunda camunda
```

## Manage a payload

To increase/decrease payload change in a file `camunda-sre-interview-master/k8s_resources/processStarter.yml`
value for env variable `N_PROCESS_STARTED` or `QUIET_TIME_S` or both

## Clean

```bash
helm delete camunda
cd camunda-sre-interview-master/
kubectl delete -f k8s_resources/
```

### Pros'n'Cons

(+) independence

(+) easy to setup

(+) no additional dependences

(-) no stabilisation window

(-) no monitoring endpoints

(-) needs addition RBAC permissions