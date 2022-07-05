#!/usr/bin/env bash
# Page explaining HID reports (the data)
# https://eleccelerator.com/tutorial-about-usb-hid-report-descriptors/
echo -ne "\0\0\x4\0\0\0\0\0" > /dev/hidg0 #press the A-button
echo -ne "\0\0\0\0\0\0\0\0" > /dev/hidg0 #release all keys
echo -ne "\0\0\x28\0\0\0\0\0" > /dev/hidg0 #press the Enter-button
echo -ne "\0\0\0\0\0\0\0\0" > /dev/hidg0 #release all keys
