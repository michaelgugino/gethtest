#!/bin/bash

pushd .
cd /tmp
wget https://go.dev/dl/go1.19.5.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xf go1.19.5.linux-amd64.tar.gz
popd
