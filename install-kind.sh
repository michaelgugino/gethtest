#!/bin/bash
pushd .

# first install kubectl
https://dl.k8s.io/v1.25.6/kubernetes-client-linux-amd64.tar.gz
tar xf kubernetes-client-linux-amd64.tar.gz
mv kubernetes/client/bin/kubectl /usr/local/bin

go install sigs.k8s.io/kind@v0.17.0

popd
