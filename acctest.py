import matplotlib.pyplot as plt
import numpy as np
import sys
import os

args = sys.argv
acc = float(args[1]) / 100
circles = int(args[2])
sliders = int(args[3])
misses = int(args[4])

objects = circles + sliders
greats = objects - misses
goods = 0
mehs = 0

newacc = (6*greats + 2*goods + mehs) / (6*objects)

if newacc > acc:
    while True:
        greats -= 1
        goods += 1
        newacc = (6*greats + 2*goods + mehs) / (6*objects)
        if newacc < acc:
            goods -= 1
            mehs += 1
            newacc = (6*greats + 2*goods + mehs) / (6*objects)
            if newacc < acc:
                mehs -= 1
                greats += 1
                if goods == 0:
                    break
                else:
                    while True:
                        if goods == 0:
                            break
                        goods -= 1
                        mehs += 1
                        newacc = (6*greats + 2*goods + mehs) / (6*objects)
                        if newacc < acc:
                            goods += 1
                            mehs -= 1
                            break

                break

print(greats, goods, mehs, misses, objects, greats+goods+mehs+misses)
