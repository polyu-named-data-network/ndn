#!/usr/bin/env bash
find -type f -name *.go -exec go fmt -x {} \;
#find -type f -name *.go -exec sed -i 's/\t/  /g' {} \;
