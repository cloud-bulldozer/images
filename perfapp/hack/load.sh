#!/bin/bash

COPIES=${1:-40}

for i in $(seq ${COPIES}); do
  echo $i
  oc process -f ../deploy/perf-app.yml IDENTIFIER=${i} | oc apply -f -
done
