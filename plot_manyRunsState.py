#!/usr/bin/python

import json, sys, math
import matplotlib 
import matplotlib.pyplot as plt


t20 = [(31, 119, 180), (174, 199, 232), (255, 127, 14), (255, 187, 120),  
             (44, 160, 44), (152, 223, 138), (214, 39, 40), (255, 152, 150),  
             (148, 103, 189), (197, 176, 213), (140, 86, 75), (196, 156, 148),  
             (227, 119, 194), (247, 182, 210), (127, 127, 127), (199, 199, 199),  
             (188, 189, 34), (219, 219, 141), (23, 190, 207), (158, 218, 229)] 

# Scale the RGB values to the [0, 1] range
for i in range(len(t20)):
    r, g, b = t20[i]  
    t20[i] = (r / 255., g / 255., b / 255.)  

# read files
dataA = {}
dataB = {}
nCol = 0
for directory in sys.argv[1:]:
	with open(directory + '/GeneticsA.json', 'r') as f:
		dataA = json.load(f)
	with open(directory + '/GeneticsB.json', 'r') as f:
		dataB = json.load(f)
	A = [point['Attraction'] for point in dataA]
	B = [point['Attraction'] for point in dataB]

	# x is A -> A
	AtoA = [p['s'] for p in A[-1:]]
	# y is B -> A
	BtoA = [p['o'] for p in B[-1:]]

	plt.scatter(AtoA, BtoA, marker='*', c=t20[nCol])

	# x is A -> B
	AtoB = [p['o'] for p in A[-1:]]
	# y is B -> B
	BtoB = [p['s'] for p in B[-1:]]

	plt.scatter(AtoB, BtoB, marker='+', c=t20[nCol])
	nCol = (nCol + 1) % len(t20)


x1,x2,y1,y2 = plt.axis()

low = min(x1, y1)
high = max(x2, y2)

plt.axis((low,high,low,high))
plt.show()