# Prow

Base image for running OCP Perf&Scale prow jobs.

## Base

- UBI 9 (`registry.access.redhat.com/ubi9/ubi:latest`)

## Included tools

| Tool | Description |
|------|-------------|
| Python 3.11 | With pip, virtualenv |
| oc / kubectl | OpenShift and Kubernetes CLIs (stable channel) |
| rosa | Red Hat OpenShift on AWS CLI |
| ocm | OpenShift Cluster Manager CLI |
| kube-burner | Kubernetes load testing tool |
| kube-burner-ocp | OCP-specific kube-burner variant |
| ocp-metadata | Cluster metadata collector from go-commons |
| orion | Regression analysis tool (v1.1.4) |
| yq | YAML processor |
| jq | JSON processor |
| gsutil | Google Cloud Storage utility |
| git | Version control |
| rsync / sshpass / openssh-clients | Remote file transfer and SSH utilities |
| ncat / iputils | Network diagnostic tools |
| gettext / util-linux / procps-ng | System utilities |

## Build

```bash
podman build -t quay.io/cloud-bulldozer/prow:latest prow/
```