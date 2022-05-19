# ein R-Gespräch

this kubernetes utility will call back a webhook url...

## Testing

```shell
# get a Kind cluster and deploy cert-manager
kind create cluster
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
# build and load the images we need for the controller
make docker-build
kind load docker-image r-gespraech-controller:v0.0.2 --name kind
# deploy our Operator
make deploy IMG=r-gespraech-controller:v0.0.2
```

## TODO

- [WIP] watch for CallbackPayload create, so that CallbackURL reconciler runs

## References

- <https://maelvls.dev/kubernetes-conditions/#are-conditions-still-used>
