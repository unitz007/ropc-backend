#!/bin/bash

git add .
git commit -m "$1"

latest=$(cat latest-image.txt)
latest=$(expr "$latest" + 1)
echo "$latest"
