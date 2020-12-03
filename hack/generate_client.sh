#!/usr/bin/env bash

printf 'package hack\nimport "k8s.io/code-generator"\n' > ./hack/hold.go
go mod vendor
retVal=$?
rm -f ./hack/hold.go
if [ $retVal -ne 0 ]; then
    exit $retVal
fi

set -e
TMP_DIR=$(mktemp -d)
mkdir -p "${TMP_DIR}"/src/github.com/SchoIsles/kruise-state-metrics/pkg/client
cp -r ./{hack,vendor} "${TMP_DIR}"/src/github.com/SchoIsles/kruise-state-metrics/

(cd "${TMP_DIR}"/src/github.com/SchoIsles/kruise-state-metrics; \
    GOPATH=${TMP_DIR} GO111MODULE=off /bin/bash vendor/k8s.io/code-generator/generate-groups.sh client \
    github.com/SchoIsles/kruise-state-metrics/pkg/client github.com/openkruise/kruise-api apps:v1alpha1 \
     -h ./hack/boilerplate.go.txt)

rm -rf ./pkg/client
mv "${TMP_DIR}"/src/github.com/SchoIsles/kruise-state-metrics/pkg/client ./pkg/client
rm -rf ${TMP_DIR}