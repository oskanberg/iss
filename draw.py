#!/usr/bin/python

import pygame
from pygame.locals import *
import time, sys, json

white = (255,255,255)
black = (0,0,0)
red = (255, 0, 0)
grey = (100,100,100)

class Pane(object):
    def __init__(self):
        pygame.init()
        pygame.display.set_caption('Simulation')
        self.screen = pygame.display.set_mode((1000,1000), 0, 32)
        self.screen.fill((white))
        pygame.display.flip()


    def drawAgent(self, agent):
    	colour = black
        render_view_circle = False
        render_view_line = False
        render_fitness = False

    	x = agent['p']['X']
    	y = agent['p']['Y']

        nx = agent['v']['X'] + x
        ny = agent['v']['Y'] + y
        
        if render_view_circle:
            pygame.draw.circle(self.screen, (grey), (int(x), int(y)), 100, 1)
        if render_fitness:
            myfont = pygame.font.SysFont("monospace", 15)
            label = myfont.render(str(agent['f']), 1, black)
            self.screen.blit(label, (x, y))

        pygame.draw.circle(self.screen, (colour), (int(x), int(y)), 5)
        if render_view_line:
            pygame.draw.line(self.screen, (colour), (x, y), (nx, ny))

    def load_and_draw(self):
    	json_data = {}
    	with open(sys.argv[1], 'r') as f:
    		json_data = json.load(f)
    	for step in json_data:
	    	for position in step:
	    		self.drawAgent(position)
	    	pygame.display.flip()
	    	time.sleep(0.001)
	    	self.screen.fill(white)	
    	# raw_input()

display = Pane()
display.load_and_draw()