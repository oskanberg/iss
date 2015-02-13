#!/usr/bin/python

import sys, json
import matplotlib.pyplot as plt

def smooth(data, window):
	print "smooth", window
	return [sum(data[i:i+window])/window for i in range(len(data)-(window-1))]


SMOOTH = 1

if len(sys.argv) == 4:
	SMOOTH = int(sys.argv[3])

plt.figure(1)
json_data = {}
with open(sys.argv[1], 'r') as f:
	json_data = json.load(f)

plt.subplot(121)
plt.title('Type A')
plt.plot(smooth(json_data, SMOOTH))
x1,x2,y1,y2 = plt.axis()
plt.axis((x1,x2,0,40))

json_data = {}
with open(sys.argv[2], 'r') as f:
	json_data = json.load(f)

plt.figure(1)
plt.subplot(122)
plt.title('Type B')
plt.plot(smooth(json_data, SMOOTH))
x1,x2,y1,y2 = plt.axis()
plt.axis((x1,x2,0,40))

	
plt.show()
