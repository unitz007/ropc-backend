#!/bin/bash

git add .
git commit -m "$1"

latest=$(cat latest-image.txt)
updated=$(expr "$latest" + 1)
echo $updated

sed -i '' 's/ropc:$latest/ropc:$updated/g' argo-cd/Deployment.yaml