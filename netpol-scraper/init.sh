#!/bin/bash

set -ex

python3.11 -m http.server &
python3.11 pod-scraper.py
