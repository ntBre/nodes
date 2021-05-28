import json

with open('pbs.json') as f:
    lines = json.load(f)

for x in lines['nodes']:
    print(lines['nodes'][x]['Mom'])
