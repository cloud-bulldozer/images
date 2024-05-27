#!/bin/bash

set -ex

python3 -m http.server &
python pod-scraper.py
