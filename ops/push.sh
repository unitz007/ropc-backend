#!/bin/bash

git add .
git commit -m "$1"

latest=$(cat latest-image.txt)
# shellcheck disable=SC2046
val=$latest + 0.1
echo "$val"
