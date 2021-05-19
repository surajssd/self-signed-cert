# Self Signed Certificate

This tool generates self signed certificates that can be used with Kubernetes webhook servers.

## Install

```
git clone https://github.com/surajssd/self-signed-cert
go install
```

Now the binary `self-signed-cert` will be built and available in you `GOBIN` directory.

## Usage

```
self-signed-cert --namespace <k8s namespace> --service-name <k8s service name>
```
