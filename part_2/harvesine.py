import json
import math

def calculateHarv(pair, earthRadius):
    lat1 = pair['y0']
    lat2 = pair['y1']
    lon1 = pair['x0']
    lon2 = pair['x1']

    dLat = lat2 - lat1
    dLon = lon2 - lon1

    a = math.sin(dLat/2)**2 + math.cos(lat1)* math.cos(lat2) * math.sin(dLon/2)**2
    c = 2*math.asin(math.sqrt(a))
    return earthRadius * c

dictionary = {
    "pairs": []
}

f = open('result.json')

data = json.load(f)

sum = 0
count = 0
for pair in data['pairs']:
    sum += calculateHarv(pair, 6372.8)
    count += 1
f.close()

print("count: ", count)
print("sum: ", sum)
