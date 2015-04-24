#!/usr/bin/python

import sys, json
import matplotlib.pyplot as plt

def smooth(data, window):
	return [sum(data[i:i+window])/window for i in range(len(data)-(window-1))]

# EINTERVAL = int(sys.argv[3])
SMOOTH = 10

json_data = {}
with open(sys.argv[1], 'r') as f:
	json_data = json.load(f)

same_data = []
other_data = []
for step in json_data:
	same_data.append(step['s'])
	other_data.append(step['o'])

same_data = smooth(same_data, SMOOTH)
other_data = smooth(other_data, SMOOTH)

plt.suptitle('Mean number of nearby prey', fontsize=15)
plt.figure(1)
	# plt.xlabel('Number of predators')
	# plt.ylabel('Number of prey')
plt.subplot(121)
plt.title('Type A')
plt.plot(same_data, 'g', label="Type A")
plt.plot(other_data, 'b', label="Type B")
plt.legend(loc='best', prop={'size':8})

maxy = max([max(same_data), max(other_data)])
miny = min([min(same_data), min(other_data)])

json_data = {}
with open(sys.argv[2], 'r') as f:
	json_data = json.load(f)

same_data = []
other_data = []
for step in json_data:
	same_data.append(step['s'])
	other_data.append(step['o'])

same_data = smooth(same_data, SMOOTH)
other_data = smooth(other_data, SMOOTH)

maxy = max([maxy, max(same_data), max(other_data)])
miny = min([miny, min(same_data), min(other_data)])

x1,x2,y1,y2 = plt.axis()
plt.axis((x1,x2,miny,maxy))

plt.figure(1)
plt.subplot(122)
plt.title('Type B')
plt.plot(other_data, 'g', label="Type A")
plt.plot(same_data, 'b', label="Type B")
x1,x2,y1,y2 = plt.axis()
plt.axis((x1,x2,miny,maxy))

plt.legend(loc='best', prop={'size':8})
plt.show()