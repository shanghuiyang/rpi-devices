#! /usr/bin/python

import serial

ser = serial.Serial("/dev/ttyAMA0",9600)
ser.write("AT+CGNSPWR=1\r\n")
ser.flushInput()

ser.write("AT+CGNSTST=1\r\n")
ser.flushInput()

print "done"
