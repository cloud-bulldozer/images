#!/bin/bash

set -x

ENGINE=podman
CL=quay.io/openshift/origin-tests:4.3

${ENGINE} run --rm -v ../deploy:/deploy:z -v ${KUBECONFIG:-~/.kube/config}:/root/.kube/config:z --rm  -it ${CL} /bin/bash -c 'KUBECONFIG=/root/.kube/config VIPERCONFIG=/deploy/clusterloader.yml openshift-tests run-test "[Feature:Performance][Serial][Slow] Load cluster should load the cluster [Suite:openshift]"'
