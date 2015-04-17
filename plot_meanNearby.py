#!/usr/bin/python

import sys, json
import matplotlib.pyplot as plt

def smooth(data, window):
	return [sum(data[i:i+window])/window for i in range(len(data)-(window-1))]


EINTERVAL = int(sys.argv[3])
SMOOTH = int(sys.argv[4])

json_data = {}
with open(sys.argv[1], 'r') as f:
	json_data = json.load(f)

same_data = []
other_data = []
for step in json_data:
	same_data.append(step['s'])
	other_data.append(step['o'])


same_avg = []
other_avg = []

for i in xrange(0, len(same_data)-1 , EINTERVAL):
	try:
		gen = [same_data[i + x] for x in xrange(EINTERVAL)]
		avg = sum(gen) / len(gen)
		same_avg.append(avg)

		gen = [other_data[i + x] for x in xrange(EINTERVAL)]
		avg = sum(gen) / len(gen)
		other_avg.append(avg)
	except Exception, e:
		print e

same_avg = smooth(same_avg, SMOOTH)
other_avg = smooth(other_avg, SMOOTH)

plt.suptitle('Mean number of nearby prey', fontsize=15)
plt.figure(1)
	# plt.xlabel('Number of predators')
	# plt.ylabel('Number of prey')
plt.subplot(121)
plt.title('Type A')
plt.plot(same_avg, 'g', label="Type A")
plt.plot(other_avg, 'b', label="Type B")
plt.legend(loc='best', prop={'size':8})

maxy = max([max(same_avg), max(other_avg)])
miny = min([min(same_avg), min(other_avg)])

json_data = {}
with open(sys.argv[2], 'r') as f:
	json_data = json.load(f)

same_data = []
other_data = []
for step in json_data:
	same_data.append(step['s'])
	other_data.append(step['o'])


same_avg = []
other_avg = []

for i in xrange(0, len(same_data)-1 , EINTERVAL):
	try:
		gen = [same_data[i + x] for x in xrange(EINTERVAL)]
		avg = sum(gen) / len(gen)
		same_avg.append(avg)

		gen = [other_data[i + x] for x in xrange(EINTERVAL)]
		avg = sum(gen) / len(gen)
		other_avg.append(avg)
	except Exception,e:
		print e

same_avg = smooth(same_avg, SMOOTH)
other_avg = smooth(other_avg, SMOOTH)

maxy = max([maxy, max(same_avg), max(other_avg)])
miny = min([miny, min(same_avg), min(other_avg)])

x1,x2,y1,y2 = plt.axis()
plt.axis((x1,x2,miny,maxy))

plt.figure(1)
plt.subplot(122)
plt.title('Type B')
plt.plot(other_avg, 'g', label="Type A")
plt.plot(same_avg, 'b', label="Type B")
x1,x2,y1,y2 = plt.axis()
plt.axis((x1,x2,miny,maxy))

plt.legend(loc='best', prop={'size':8})
plt.show()