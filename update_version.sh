#!/bin/sh

go mod tidy
git add .
git status

# commit
current_dir=$(basename "$PWD")
git commit -m "--feat: update go package [$current_dir]"
git push