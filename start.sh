#!/bin/bash

raspivid -cd MJPEG -w 640 -h 360 -fps 10 -t 0 -n -o - | go run main.go