# Autoscale

This is autoscaler based on client-go library and works by helm deployment

## Build

```bash
cd camunda-autoscaler/default/src
docker build -t avguston/camunda:stable
docker push avguston/camunda:stable
```


## Install

```bash
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