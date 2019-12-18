import sys
import numpy as np
import matplotlib.pyplot as plt


iter_colors = iter(["#0ED4F9","#83BAFA","#BF9CE1","#E6759F","#D66759"])


if __name__ == "__main__":

    a = np.transpose(np.loadtxt("cache/graph.txt"))
    
    times, IPs_raw, IPs, miss_probs = a[0], a[1], a[2], a[3]

    fig, axarr = plt.subplots(2, sharex=True, figsize=[12,6])
    
    axarr[0].plot(times, IPs, '.', alpha=0.8, markersize=4)
    axarr[0].vlines(times, IPs_raw, IPs, colors=(1.0,0.5,0.5,0.8), linewidths=1)
    
    axarr[0].set_ylabel("Index of Performance (bits/s)")

    axarr[1].plot(times, miss_probs, '.', alpha=0.8, markersize=4)
    axarr[1].set_xlabel("Time (s)")
    axarr[1].set_ylabel("Miss Probability")
    plt.savefig("cache/graph.png")
    plt.show()