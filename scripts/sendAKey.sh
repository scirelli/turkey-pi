#!/usr/bin/env bash
# Page explaining HID reports (the data)
# https://eleccelerator.com/tutorial-about-usb-hid-report-descriptors/
for i in {1..50}; do
	echo "$i"
	echo -ne "\0\0\x4\0\0\0\0\0" > /dev/hidg0 #press the A-button
	sleep 0.1
	echo -ne "\0\0\0\0\0\0\0\0" > /dev/hidg0 #release all keys
	sleep 0.1
done
sleep 0.1
echo -ne "\0\0\x28\0\0\0\0\0" > /dev/hidg0 #press the Enter-button
sleep 0.1
echo -ne "\0\0\0\0\0\0\0\0" > /dev/hidg0 #release all keys
