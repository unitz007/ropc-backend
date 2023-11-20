#!/bin/zsh

git add .
git commit -m "$1"

latest=$(cat latest-image.txt)+1
