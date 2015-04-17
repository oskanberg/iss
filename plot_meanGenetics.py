#!/usr/bin/python

import sys, json
import matplotlib.pyplot as plt


YLIM = 100

f1_data = {}
with open(sys.argv[1], 'r') as f:
	f1_data = json.load(f)


f2_data = {}
if len(sys.argv) > 2:
	with open(sys.argv[2], 'r') as f:
		f2_data = json.load(f)


attraction_same = []
attraction_other = []

orientation_same = []
orientation_other = []
for step in f1_data:
	attraction_same.append(step['Attraction']['s'])
	attraction_other.append(step['Attraction']['o'])

	orientation_same.append(step['Orientation']['s'])
	orientation_other.append(step['Orientation']['o'])

plt.suptitle('Mean genetics', fontsize=15)

print sys.argv
if '-alt' in sys.argv:
	plt.figure(1)
	plt.subplot(221)
	plt.title('same')
	plt.plot(orientation_same, 'b', label="Orientation")
	plt.plot(attraction_same, 'k', label="Attraction")
	x1,x2,y1,y2 = plt.axis()
	plt.axis((x1,x2,0,YLIM))
	plt.legend(loc='best', prop={'size':8})

	plt.subplot(222)
	plt.title('other')
	plt.plot(orientation_other, 'b', label="Orientation")
	plt.plot(attraction_other, 'k', label="Attraction") 
	x1,x2,y1,y2 = plt.axis()
	plt.axis((x1,x2,0,YLIM))
	plt.legend(loc='best', prop={'size':8})

	if len(sys.argv) > 2:
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
		plt.title('same')
		plt.plot(orientation_same, 'b', label="Orientation")
		plt.plot(attraction_same, 'k', label="Attraction")
		x1,x2,y1,y2 = plt.axis()
		plt.axis((x1,x2,0,YLIM))
		plt.legend(loc='best', prop={'size':8})

		plt.subplot(224)
		plt.title('other')
		plt.plot(orientation_other, 'b', label="Orientation")
		plt.plot(attraction_other, 'k', label="Attraction")
		x1,x2,y1,y2 = plt.axis()
		plt.axis((x1,x2,0,YLIM))
		plt.legend(loc='best', prop={'size':8})
else:
	plt.figure(1)
	plt.subplot(221)
	plt.title('attraction')
	plt.plot(attraction_same, 'g', label='Same subspecies')
	plt.plot(attraction_other, 'r', label='Other subspecies')
	plt.legend(loc='best', prop={'size':8})
	
	plt.subplot(222)
	plt.title('orientation')
	plt.plot(orientation_same, 'g', label='Same subspecies')
	plt.plot(orientation_other, 'r', label='Other subspecies')
	plt.legend(loc='best', prop={'size':8})

	if len(sys.argv) > 2:
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
		plt.title('attraction')
		plt.plot(attraction_same, 'g')
		plt.plot(attraction_other, 'r')

		plt.subplot(224)
		plt.title('orientation')
		plt.plot(orientation_same, 'g')
		plt.plot(orientation_other, 'r')

plt.show()