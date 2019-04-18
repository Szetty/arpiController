#!/bin/bash

mjpg_streamer -i "input_raspicam.so -vf -hf" -o "output_http.so -w ./www -p 8000" & ./arpiController -withMotors