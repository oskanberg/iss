#!/usr/bin/python

import pygame
from pygame.locals import *
import time, sys, json, math

SCALE = 2

white = (255,255,255)
black = (0,0,0)
red = (255, 0, 0)
green = (0, 255, 0)
blue = (0, 0, 255)
grey = (100,100,100)

class Pane(object):
    def __init__(self):
        pygame.init()
        pygame.display.set_caption('Simulation')
        self.screen = pygame.display.set_mode((400*SCALE,400*SCALE), 0, 32)
        self.screen.fill((white))
        pygame.display.flip()


    def drawAgent(self, agent):
    	colour = green
        render_orientation = False
        render_attraction = False
        render_view_line = True
        render_fitness = False
        render_family = False
        render_ap = False

        if agent['t'] == 2:
            colour = red
        elif agent['t'] == 1:
            colour = blue
        elif agent['t'] == 0:
            colour = green

    	x = agent['p']['X'] * SCALE
    	y = agent['p']['Y'] * SCALE

        nx = agent['v']['X'] + x
        ny = agent['v']['Y'] + y
        
        if render_orientation:
            r = 0
            try:
                r = agent['g']['os']['o']
                if r > 0:
                    r = int(math.sqrt(r))
                if r < 1:
                    r = 1   
                pygame.draw.circle(self.screen, (red), (int(x), int(y)), r * SCALE, 1)
            except Exception, e:
                print e, agent
            
            try:
                r = agent['g']['ss']['o']
                if r > 0:
                    r = int(math.sqrt(r))
                if r < 1:
                    r = 1  
                pygame.draw.circle(self.screen, (green), (int(x), int(y)), r * SCALE, 1)
            except ValueError, e:
                print agent['g']['ss']['o']
            except TypeError, e:
                pass
                print agent

        if render_attraction:
            try:
                r = agent['g']['os']['a']
                if r > 0:
                    r = int(math.sqrt(r))
                if r < 1:
                    r = 1   
                pygame.draw.circle(self.screen, (red), (int(x), int(y)),r * SCALE, 1)
            except ValueError, e:
                print agent['g']['os']['a'], e
            except TypeError, e:
                pass
                print agent

            try:
                r = agent['g']['ss']['a']
                if r > 0:
                    r = int(math.sqrt(r))
                if r < 1:
                    r = 1            
                pygame.draw.circle(self.screen, (green), (int(x), int(y)), r * SCALE, 1)
            except ValueError, e:
                print agent['g']['ss']['a'], e
            except TypeError, e:
                pass
                print agent
            
        if render_fitness:
            myfont = pygame.font.SysFont("monospace", 15)
            label = myfont.render(str(agent['f']), 1, black)
            self.screen.blit(label, (x, y))
        if render_family:
            myfont = pygame.font.SysFont("monospace", 15)
            label = myfont.render(str(agent['t']), 1, black)
            self.screen.blit(label, (x, y))
        if render_ap:
            myfont = pygame.font.SysFont("monospace", 15)
            label = myfont.render(str(agent['g']['pr']), 1, black)
            self.screen.blit(label, (x, y))


        try:
            pygame.draw.circle(self.screen, (colour), (int(x), int(y)), 2)
        except Exception, e:
            print x, y
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
	    	time.sleep(0.02)
	    	self.screen.fill(white)	
    	# raw_input()

display = Pane()
display.load_and_draw()