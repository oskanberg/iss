
from numpy import random, sqrt, log, sin, cos, pi
from pylab import *
from numpy import array, transpose
import json

# # read files
# f1_data = {}
# with open('output/DetailGeneticsA.json', 'r') as f:
# 	f1_data = json.load(f)


a = array([[2,3,4],[1,2,3],[1,2,3],[1,2,3]])
boxplot(transpose(a), positions=[1,2,4,5])
# boxplot(transpose(a), positions=[3,4])
show()