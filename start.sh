#!/bin/bash

raspivid -cd MJPEG -w 640 -h 360 -fps 1 -t 0 -n -vf -hf -o - | go run main.go
