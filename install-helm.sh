#!/bin/bash

pushd .
cd /tmp
wget https://get.helm.sh/helm-v3.11.0-linux-amd64.tar.gz
tar xf helm-v3.11.0-linux-amd64.tar.gz
mv linux-amd64/helm /usr/local/bin/
popd
