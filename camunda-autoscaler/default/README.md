# Autoscale

This is autoscaler based on client-go library and works by helm deployment

## Build

```bash
cd camunda-autoscaler/default
docker build -t avguston/camunda:stable .
docker push avguston/camunda:stable
```


## Install

```bash
cd camunda-sre-interview-master/
kubectl apply -f k8s_resources
cd camunda-autoscaler/default
helm install camunda camunda
```

## Clean

```bash
helm delete camunda
kubectl delete deployment camunda-deployment
kubectl delete deployment camunda-process-starter
```

### Pros'n'Cons

+ independence 
+ easy to setup
+ no additional dependences

- no stabilisation window
- no monitoring endpoints
- needs addition RBAC permissions