"""Diff compare-states.mjs captures (pixels at threshold 2, plus computed styles).

    python3 -I scripts/compare-states.py /tmp/cmp-states
"""
import json, sys, glob, os
from PIL import Image, ImageChops
out=sys.argv[1]; bad=0
for theme in ['light','dark']:
    m=json.load(open(f'{out}/{theme}-mock.json')); w=json.load(open(f'{out}/{theme}-web.json'))
    for name in m:
        if m[name]=='MISSING' or w.get(name)=='MISSING': print(f'{theme} {name}: MISSING mock={m[name]=="MISSING"} web={w.get(name)=="MISSING"}'); bad+=1; continue
        for state in m[name]:
            a=f'{out}/{theme}-{name}-{state}-mock.png'; b=f'{out}/{theme}-{name}-{state}-web.png'
            A=Image.open(a).convert('RGB'); B=Image.open(b).convert('RGB')
            if A.size!=B.size: print(f'{theme} {name:13} {state:7} SIZE {A.size} vs {B.size}'); bad+=1; continue
            d=ImageChops.difference(A,B).convert('L').point(lambda v:255 if v>2 else 0)
            n=sum(1 for v in d.get_flattened_data() if v)
            sd={k:(m[name][state][k],w[name][state][k]) for k in m[name][state] if m[name][state][k]!=w[name][state][k]}
            if n or sd: bad+=1; print(f'{theme} {name:13} {state:7} px={n} {json.dumps(sd)[:230]}')
print('differences:',bad)
