# ein R-Gespräch

this kubernetes utility will call back a webhook url...

## Testing

### locally on a Kind cluster

```shell
# get a Kind cluster and deploy cert-manager
kind create cluster --config hack/kind-config.yaml
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.5.0/aio/deploy/recommended.yaml
kubectl apply -f hack/dashboard-adminuser.yaml
export T=$(kubectl -n kubernetes-dashboard create token admin-user)
# build and load the images we need for the controller
make docker-build IMG=r-gespraech-controller:v0.0.2
kind load docker-image r-gespraech-controller:v0.0.2 --name kind
# deploy our Operator
make deploy IMG=r-gespraech-controller:v0.0.2
```

### in an envtest

If you want to see some output enter the `controllers/` directory and run the test: `go test ./... -ginkgo.v -v`
or just `make test` on the root of the repository.

## TODO

- [WIP] double-check if this is a cluster scoped operator
- [WIP] watch for CallbackPayload create, so that CallbackURL reconciler runs

## References

- <https://maelvls.dev/kubernetes-conditions/#are-conditions-still-used>
- <https://itnext.io/big-change-in-k8s-1-24-about-serviceaccounts-and-their-secrets-4b909a4af4e0>
