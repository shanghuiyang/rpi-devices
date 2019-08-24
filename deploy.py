#! /usr/bin/python

import sys
import os

file = 'main.go'
bin = 'devices.pi'
for i in range(1, len(sys.argv)):
    if sys.argv[i] == "-test":
        bin = 'test.pi'
    if sys.argv[i] == '-f':
        file = sys.argv[i+1]
    if sys.argv[i] == '-o':
        bin = sys.argv[i+1]

print('{} --> {}'.format(file, bin))
print("building for pi...")
res = os.system('CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o {} {}'.format(bin, file))
if res != 0:
    print('\033[1;31m[failed]\033[0m')  # highlight in red
    exit(-1)

print('\033[1;32m[success]\033[0m') # highlight in green

print("deploying {}...".format(bin))
res = os.system('scp {} pi@192.168.31.57:/home/pi'.format(bin))
if res == 0:
    print('\033[1;32m[success]\033[0m') # highlight in green
    exit(0)
else:
    print('\033[1;31m[failed]\033[0m')  # highlight in red
    exit(-1)
