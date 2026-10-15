#!/bin/bash

set -e

# Create etcd write dir if it doesn't exist
mkdir -p /var/lib/etcd

# Run fio
echo "---------------------------------------------------------------- Running fio ---------------------------------------------------------------------------"
fio --rw=write --ioengine=sync --fdatasync=1 --directory=/var/lib/etcd --size=100m --bs=8000 --name=etcd_perf --output-format=json --runtime=60 --time_based=1 | tee /tmp/fio.out
echo "--------------------------------------------------------------------------------------------------------------------------------------------------------"

# Scrape the fio output for p99 and p99.9 of fsync in ns
fsync_p99=$(cat /tmp/fio.out | jq '.jobs[0].sync.lat_ns.percentile["99.000000"]')
fsync_p99_9=$(cat /tmp/fio.out | jq '.jobs[0].sync.lat_ns.percentile["99.900000"]')

echo "INFO: 99th percentile of fsync is $fsync_p99 ns ($(awk "BEGIN {printf \"%.1f\", $fsync_p99/1000000}") ms)"
echo "INFO: 99.9th percentile of fsync is $fsync_p99_9 ns ($(awk "BEGIN {printf \"%.1f\", $fsync_p99_9/1000000}") ms)"

# Threshold: 10ms = 10,000,000 ns
THRESHOLD=10000000

# Check both p99 and p99.9 against recommended thresholds
# Per Red Hat KB 4885641, BOTH p99 and p99.9 must be under 10ms
p99_pass=false
p99_9_pass=false

if [[ $fsync_p99 -lt $THRESHOLD ]]; then
  p99_pass=true
  echo "INFO: p99 fsync latency is within the recommended threshold (< 10 ms)"
else
  echo "ERROR: p99 fsync latency ($fsync_p99 ns = $(awk "BEGIN {printf \"%.1f\", $fsync_p99/1000000}") ms) exceeds the recommended threshold of 10 ms"
fi

if [[ $fsync_p99_9 -lt $THRESHOLD ]]; then
  p99_9_pass=true
  echo "INFO: p99.9 fsync latency is within the recommended threshold (< 10 ms)"
else
  echo "ERROR: p99.9 fsync latency ($fsync_p99_9 ns = $(awk "BEGIN {printf \"%.1f\", $fsync_p99_9/1000000}") ms) exceeds the recommended threshold of 10 ms"
fi

echo "--------------------------------------------------------------------------------------------------------------------------------------------------------"

# Final verdict: BOTH must pass
if [[ "$p99_pass" == "true" && "$p99_9_pass" == "true" ]]; then
  echo "PASS: Both p99 and p99.9 fsync latencies are within recommended thresholds"
  echo "INFO: The disk can be used to host etcd"
  exit 0
else
  echo "FAIL: Disk does NOT meet Red Hat requirements for hosting etcd"
  echo "INFO: Per Red Hat KB 4885641, both p99 and p99.9 fsync latencies must be under 10 ms"
  echo "INFO: Faster disks are recommended to host etcd for better performance"
  exit 1
fi
