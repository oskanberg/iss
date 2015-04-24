#!/usr/bin/python

import json, sys
import matplotlib 
import matplotlib.pyplot as plt

# read files
f1_data = {}
with open(sys.argv[1], 'r') as f:
	f1_data = json.load(f)


f2_data = {}
with open(sys.argv[2], 'r') as f:
	f2_data = json.load(f)


# Parse into series
attraction_same = []
attraction_other = []
orientation_same = []
orientation_other = []
for step in f1_data:
	attraction_same.append(step['Attraction']['s'])
	attraction_other.append(step['Attraction']['o'])
	orientation_same.append(step['Orientation']['s'])
	orientation_other.append(step['Orientation']['o'])


# plot
fig, axes = plt.subplots(nrows=2, ncols=2)

plt.subplot(221)
plt.ylim((0,100))
plt.title('Same subspecies')
plt.ylabel('Type A', rotation=90, size='large')
plt.plot(orientation_same, 'r', label="Orientation")
plt.plot(attraction_same, 'k', label="Attraction")
plt.legend(loc='best', prop={'size':8})

plt.subplot(222)
plt.ylim((0,100))
plt.title('Other subspecies')
plt.plot(orientation_other, 'r', label="Orientation")
plt.plot(attraction_other, 'k', label="Attraction")
plt.legend(loc='best', prop={'size':8})

attraction_same = []
attraction_other = []
orientation_same = []
orientation_other = []
for step in f2_data:
	attraction_same.append(step['Attraction']['s'])
	attraction_other.append(step['Attraction']['o'])
	orientation_same.append(step['Orientation']['s'])
	orientation_other.append(step['Orientation']['o'])

plt.subplot(223)
plt.ylim((0,100))
plt.ylabel('Type B', rotation=90, size='large')
plt.plot(orientation_same, 'r', label="Orientation")
plt.plot(attraction_same, 'k', label="Attraction")
plt.legend(loc='best', prop={'size':8})

plt.subplot(224)
plt.ylim((0,100))
plt.plot(orientation_other, 'r', label="Orientation")
plt.plot(attraction_other, 'k', label="Attraction")
plt.legend(loc='best', prop={'size':8})

# fig.tight_layout()
plt.show()
