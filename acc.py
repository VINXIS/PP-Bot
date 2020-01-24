import matplotlib.pyplot as plt
import numpy as np
import sys
import os

args = sys.argv
mapInfo = args[1]
difference = float(args[4])
ticks = float(args[5])

accVals = []
aimVals = []
tapVals = []
accppVals = []
ppVals = []
if os.path.exists(mapInfo+'ppVals.txt'):
    for t in open(mapInfo+'ppVals.txt').read().split('\n'):
        if '(' in t:
            acc, aim, tap, accpp, pp = t.strip('()').split(',')
            accVals.append(float(acc))
            aimVals.append(float(aim))
            tapVals.append(float(tap))
            accppVals.append(float(accpp))
            ppVals.append(float(pp))
fig = plt.figure(figsize=[12, 4])
plt.plot(accVals, aimVals, label='aim')
plt.plot(accVals, tapVals, label='tap')
plt.plot(accVals, accppVals, label='acc')
plt.plot(accVals, ppVals, label='PP')

plt.title(mapInfo)
plt.xlabel('acc')
plt.ylabel('PP')
plt.legend()
plt.xlim(float(args[2]), float(args[3]))
plt.xticks(np.arange(float(args[2]), float(args[3]), ticks))
plt.grid(b=True, which='both')
plt.savefig(mapInfo + '.png', bbox_inches='tight')