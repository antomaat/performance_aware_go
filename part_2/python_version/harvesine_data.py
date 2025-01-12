import json
import random

dictionary = {
    "pairs": []
}

for x in range(100):
    pair = {
        "x0": random.randint(-180, 180),
        "y0": random.randint(-180, 180),
        "x1": random.randint(-180, 180),
        "y1": random.randint(-180, 180)
    }
    dictionary["pairs"].append(pair)

result = json.dumps(dictionary, indent=4)

with open("result.json", "w") as outfile:
    outfile.write(result)
    outfile.write("\n")

