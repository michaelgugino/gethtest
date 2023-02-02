#!/bin/bash

./install-go.sh
export PATH=$PATH:/usr/local/go/bin

./install-docker.sh
./install-helm.sh
./install-kubebuilder.sh
./install-kind.sh

echo "run 'export PATH=\$PATH:/usr/local/go/bin'"
echo "to add go to your path"
