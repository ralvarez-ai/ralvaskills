"""Per-section pixel diff of compare.mjs captures.

    python3 -I scripts/compare-diff.py /tmp/cmp 2 1440-light,390-dark

Threshold 2 (levels of 255): a threshold of 24 once hid a tinted section.
"""
import json, sys
from PIL import Image, ImageChops
out=sys.argv[1]; thr=int(sys.argv[2])
for cfg in sys.argv[3].split(','):
    A=Image.open(f'{out}/{cfg}-mock.png').convert('RGB'); B=Image.open(f'{out}/{cfg}-web.png').convert('RGB')
    secs=json.load(open(f'{out}/{cfg}-mock.json'))['sections']
    h=min(A.height,B.height)
    d=ImageChops.difference(A.crop((0,0,A.width,h)),B.crop((0,0,B.width,h))).convert('L').point(lambda v:255 if v>thr else 0)
    tot=sum(1 for v in d.get_flattened_data() if v)
    worst=[]
    for s in secs:
        c=d.crop((0,s['y'],d.width,min(h,s['y']+s['h']))); n=sum(1 for v in c.get_flattened_data() if v); worst.append((n,s['id']))
    worst.sort(reverse=True)
    print(f"{cfg:11} sizes {A.size}=={B.size}  diff>{thr}: {tot} px ({100*tot/(A.width*h):.4f}%)  worst: {worst[:3]}")
