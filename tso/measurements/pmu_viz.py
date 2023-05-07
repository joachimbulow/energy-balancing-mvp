import json
import matplotlib.pyplot as plt

with open('pmu_new.json', 'r') as file:
    data = json.load(file)

locations = set([d['location'] for d in data])

for location in locations:
    x = [d['timestamp'] for d in data if d['location'] == location]
    y = [d['frequency'] for d in data if d['location'] == location]
    plt.plot(x, y, label=location)

plt.xlabel('Timestamp')
plt.ylabel('Frequency')
plt.title('Frequency by Location')
plt.legend()
plt.show()