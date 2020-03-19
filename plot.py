import matplotlib.pyplot as plt
import numpy as np
import sys
import os

args = sys.argv
skill = args[1]
beatmapid = args[2]
start = round(float(args[3]))
end = round(float(args[4]))
difference = float(args[5])
mapinfo = args[6]
mods = args[7]
version = args[8]
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
    if os.path.exists("cache/graph_" + beatmapid + "_" + mods + ".txt"):
        a = np.transpose(np.loadtxt("cache/graph_" + beatmapid + "_" + mods + ".txt"))
    elif os.path.exists("cache/graph__" + mods + ".txt"):
        a = np.transpose(np.loadtxt("cache/graph__" + mods + ".txt"))

    times, IPs_raw, IPs, miss_probs = a[0], a[1], a[2], a[3]

    fig, axarr = plt.subplots(2, sharex=True, figsize=[48,6])
    
    axarr[0].plot(times, IPs, '.', alpha=0.8)
    axarr[0].vlines(times, IPs_raw, IPs, colors=(1.0,0.5,0.5,0.8), linewidths=1)
    
    axarr[0].set_ylabel("Index of Performance (bits/s)")

    axarr[1].plot(times, miss_probs, '.', alpha=0.8)
    axarr[1].set_xlabel("Time (s)")
    axarr[1].set_ylabel("Miss Probability")
    plt.title(mapinfo)
elif version == "tap":
    if os.path.exists("cache/graph_" + beatmapid + "_" + mods + "_tap.txt"):
        a = np.transpose(np.loadtxt("cache/graph_" + beatmapid + "_" + mods + "_tap.txt"))
    elif os.path.exists("cache/graph__" + mods + "_tap.txt"):
        a = np.transpose(np.loadtxt("cache/graph__" + mods + "_tap.txt"))

    times, totalStrains1, totalStrains2, totalStrains3, totalStrains4, strains = a[0], a[1], a[2], a[3], a[4], a[5]

    fig, axarr = plt.subplots(5, sharex=True, figsize=[48,15])

    axarr[0].plot(times, strains, '.', alpha=0.8)

    axarr[0].set_ylabel("Specific note strains")

    axarr[1].plot(times, totalStrains1, '.-')
    axarr[1].set_xlabel("Time (s)")
    axarr[1].set_ylabel("Total note strains 9.97418")

    axarr[2].plot(times, totalStrains2, '.-')
    axarr[2].set_xlabel("Time (s)")
    axarr[2].set_ylabel("Total note strains 1.82212")
    
    axarr[3].plot(times, totalStrains3, '.-')
    axarr[3].set_xlabel("Time (s)")
    axarr[3].set_ylabel("Total note strains 0.332871")

    axarr[4].plot(times, totalStrains4, '.-')
    axarr[4].set_xlabel("Time (s)")
    axarr[4].set_ylabel("Total note strains 0.0608101")

    plt.title(mapinfo + " - Finger Control")
elif version == "finger":
    if os.path.exists("cache/graph_" + beatmapid + "_" + mods + "_finger.txt"):
        a = np.transpose(np.loadtxt("cache/graph_" + beatmapid + "_" + mods + "_finger.txt"))
    elif os.path.exists("cache/graph__" + mods + "_finger.txt"):
        a = np.transpose(np.loadtxt("cache/graph__" + mods + "_finger.txt"))

    times, totalStrains, strains = a[0], a[1], a[2]

    fig, axarr = plt.subplots(2, sharex=True, figsize=[48,6])

    axarr[0].plot(times, strains, '.', alpha=0.8)

    axarr[0].set_ylabel("Specific note strains")

    axarr[1].plot(times, totalStrains, '.-')
    axarr[1].set_xlabel("Time (s)")
    axarr[1].set_ylabel("Total note strains")
    plt.title(mapinfo + " - Finger Control")

plt.xlim(start, end)
plt.xticks(np.arange(start, end, ticks))
plt.minorticks_on()
plt.grid(b=True, which='both')
plt.savefig(beatmapid + '.png', bbox_inches='tight')