#!/bin/bash

git add .
git commit -m "$1"

latest=$(cat latest-image.txt)
# shellcheck disable=SC2046
val=`expr $latest + 1.0`
echo "$val"
