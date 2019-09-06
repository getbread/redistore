#!/bin/bash
set -euo pipefail

echo "build image"
docker build .
