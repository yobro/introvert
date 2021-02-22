#!/usr/bin/env bash
set -e
set -u

cd web/app
echo "building web app"

PUBLIC_URL=. yarn build
rm -rf ../static/app
mv dist ../static/app