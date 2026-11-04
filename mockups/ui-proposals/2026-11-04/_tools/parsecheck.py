"""HTML parser check: any start tag outside HTML/SVG means copy was swallowed
as markup (a raw <token> or <a@b.c>). Also flags leaked template braces."""
import re
import sys
from html.parser import HTMLParser
from pathlib import Path

KNOWN = set("""html head body meta title style script link main header footer nav section article aside figure figcaption
h1 h2 h3 h4 h5 h6 p a span b i em strong small code pre kbd ul ol li dl dt dd div table caption thead tbody tfoot tr th td
button details summary label input select option svg text g rect line path circle polyline tspan title desc br hr img time abbr fieldset legend output blockquote""".split())


class P(HTMLParser):
    def __init__(self):
        super().__init__()
        self.bad, self.stack = [], []
    def handle_starttag(self, tag, attrs):
        if tag not in KNOWN:
            self.bad.append((self.getpos(), tag))


bad = 0
for f in sys.argv[1:]:
    src = Path(f).read_text()
    p = P(); p.feed(src)
    body = re.sub(r"<(script|style)[^>]*>.*?</\1>", "", src, flags=re.S)
    leaks = re.findall(r"\{[a-z_]+\}", body)
    print(("FAIL" if p.bad or leaks else "PASS"), Path(f).name, "unknown tags:", p.bad[:5], "leaks:", leaks[:5])
    bad += bool(p.bad or leaks)
sys.exit(1 if bad else 0)
