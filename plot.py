import matplotlib.pyplot as plt
import numpy as np
import sys
import os

args = sys.argv
skill = args[1]
beatmapid = args[2]
difference = float(args[5])
mapinfo = args[6]
version = args[7]
ticks = 0.1
if difference > 120:
    ticks = 20
elif difference > 60:
    ticks = 10
elif difference > 30:
    ticks = 5
elif difference > 10:
    ticks = 2
elif difference > 5:
    ticks = 1
elif difference > 2:
    ticks = 0.5

if version == "joz":
    val = []
    if os.path.exists(beatmapid+skill+'.txt'):
        for t in open(beatmapid+skill+'.txt').read().split('\n'):
            if '(' in t:
                a, b = t.strip('()').split(',')
                val.append((int(a) / 1000, float(b)))
    elif os.path.exists(skill+'.txt'):
        for s in open(skill+'.txt').read().split('\n'):
            if '(' in s:
                a, b = s.strip('()').split(',')
                val.append((int(a) / 1000, float(b)))
    fig = plt.figure(figsize=[48, 6])
    plt.plot(*zip(*[(elem1, elem2) for elem1, elem2 in val]))

    plt.title(mapinfo + " - " + skill)
    plt.xlabel('seconds')
    plt.ylabel('strain')
elif version == "delta":
    if os.path.exists("cache/graph_" + beatmapid + ".txt"):
        a = np.transpose(np.loadtxt("cache/graph_" + beatmapid + ".txt"))
    elif os.path.exists("cache/graph_.txt"):
        a = np.transpose(np.loadtxt("cache/graph_.txt"))

    times, IPs_raw, IPs, miss_probs = a[0], a[1], a[2], a[3]

    fig, axarr = plt.subplots(2, sharex=True, figsize=[48,6])
    
    axarr[0].plot(times, IPs, '.', alpha=0.8)
    axarr[0].vlines(times, IPs_raw, IPs, colors=(1.0,0.5,0.5,0.8), linewidths=1)
    
    axarr[0].set_ylabel("Index of Performance (bits/s)")

    axarr[1].plot(times, miss_probs, '.', alpha=0.8)
    axarr[1].set_xlabel("Time (s)")
    axarr[1].set_ylabel("Miss Probability")
    plt.title(mapinfo)

plt.xlim(round(float(args[3])), round(float(args[4])))
plt.xticks(np.arange(round(float(args[3])), round(float(args[4])), ticks))
plt.minorticks_on()
plt.grid(b=True, which='both')
plt.savefig(beatmapid + '.png', bbox_inches='tight')