#!/bin/bash

BASEDIR=$( dirname "${BASH_SOURCE[0]}" )
# Grafana requires world readership. Setting umask only affects the local shell.
umask 0002
for file in `ls ${BASEDIR}/*.dashboard.py` `ls ${BASEDIR}/*.jsonnet`; do
  prefix=${file%%.dashboard.py}
  prefix=${prefix%%.jsonnet}
  echo "Generating: ${prefix}.json"
  if [[ $file == *.dashboard.py ]] ; then
      python3 /usr/local/bin/generate-dashboard \
          -o ${prefix}.json ${file}
  else
      jsonnet -J /go/src/github.com/grafana/grafonnet-lib \
          -o ${prefix}.json ${file}
  fi
done
