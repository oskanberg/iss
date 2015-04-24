#!/usr/bin/python

import json, sys, math
import matplotlib 
import matplotlib.pyplot as plt

# read files
data = {}
with open(sys.argv[1], 'r') as f:
	data = json.load(f)

typeA = [ agent for agent in data[0] if agent['t'] == 0 and agent['g'] != None ]
typeB = [ agent for agent in data[0] if agent['t'] == 1 and agent['g'] != None ]

print len(data[0]), len(typeA), len(typeB)

# x is A
aX = [math.sqrt(agent['g']['ss']['a']) for agent in typeA]
aY = [math.sqrt(agent['g']['os']['a']) for agent in typeA]
# x is A
bX = [math.sqrt(agent['g']['os']['a']) for agent in typeB]
bY = [math.sqrt(agent['g']['ss']['a']) for agent in typeB]

plt.scatter(aX, aY, c="g")
plt.scatter(bX, bY, c="b")

x1,x2,y1,y2 = plt.axis()
low = min(x1, y1)
high = max(x2, y2)
plt.axis((low,high,low,high))

plt.show()