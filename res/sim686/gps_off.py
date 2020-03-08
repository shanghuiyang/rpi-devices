#! /usr/bin/python

import serial

ser = serial.Serial("/dev/ttyAMA0",9600)
ser.close()

print "done"
