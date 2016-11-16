#!/usr/bin/env bash
env GOOS=linux GOARCH=arm go build -x -o ../klingenplayer klingenplayer.go