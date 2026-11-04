// node verify.js direction-1.html [...]
// axe-core at 390x844 and 1440x800 x light/dark x motion/reduced = 8 runs per page,
// plus no horizontal scroll, images blocked, and JS off. Exit 1 on any failure.
const fs = require('fs');
const path = require('path');
const { DIR, launch, open, settle, axe, overflow, unrevealed } = require('./lib');

(async () => {
  const files = process.argv.slice(2);
  const browser = await launch();
  const report = {};
  let bad = 0;
  for (const file of files) {
    report[file] = [];
    for (const width of [390, 1440]) for (const theme of ['light', 'dark']) for (const reduced of [false, true]) {
      const { ctx, page, errors } = await open(browser, file, { width, theme, reduced });
      await settle(page, { reduced });
      const v = await axe(page);
      const o = await overflow(page);
      const hid = await unrevealed(page);
      if (hid) errors.push(`${hid} reveal(s) never fired`);
      const name = `${width} ${theme} ${reduced ? 'reduced' : 'motion'}`;
      const ok = !v.length && !o.overflow && !errors.length;
      if (!ok) bad++;
      report[file].push({ run: name, violations: v, overflow: o, errors });
      console.log(`${ok ? 'PASS' : 'FAIL'} ${file} ${name}` + (v.length ? ' axe:' + v.map((x) => `${x.id}(${x.n})`).join(',') : '') + (o.overflow ? ` overflow:${o.sw}>${o.w} ${o.culprits.join(' | ')}` : '') + (errors.length ? ' errors:' + errors.join(' | ') : ''));
      for (const x of v) for (const t of x.targets) console.log('     ', x.id, t);
      await ctx.close();
    }
    // JS off: everything readable, nothing left hidden by motion rules
    for (const width of [390, 1440]) {
      const { ctx, page } = await open(browser, file, { width, theme: 'light', reduced: false, js: false });
      const r = await page.evaluate(() => {
        const h1 = document.querySelector('h1');
        const hidden = [...document.querySelectorAll('main *')].filter((e) => { const s = getComputedStyle(e); return (s.opacity === '0' || s.visibility === 'hidden') && e.textContent.trim() && !e.closest('[hidden],.sr'); }).map((e) => e.tagName + '.' + e.className).slice(0, 5);
        return { h1: h1 && h1.innerText.trim(), hidden, sw: document.documentElement.scrollWidth, w: document.documentElement.clientWidth };
      });
      const ok = r.h1 && !r.hidden.length && r.sw <= r.w;
      if (!ok) bad++;
      report[file].push({ run: `${width} js-off`, ...r });
      console.log(`${ok ? 'PASS' : 'FAIL'} ${file} ${width} js-off h1="${(r.h1 || '').slice(0, 40)}" hidden=${r.hidden.length} ${r.hidden.join(',')}`);
      await ctx.close();
    }
    // images blocked: no <img> is load-bearing; the page must still read
    {
      const { ctx, page, errors } = await open(browser, file, { width: 390, theme: 'dark', reduced: true, blockImages: true });
      const r = await page.evaluate(() => ({ imgs: document.images.length, noAlt: [...document.images].filter((i) => !i.hasAttribute('alt')).length, text: document.querySelector('main').innerText.length }));
      const ok = r.noAlt === 0 && r.text > 500 && !errors.length;
      if (!ok) bad++;
      report[file].push({ run: 'images-blocked', ...r });
      console.log(`${ok ? 'PASS' : 'FAIL'} ${file} images-blocked imgs=${r.imgs} text=${r.text}`);
      await ctx.close();
    }
    const size = fs.statSync(path.join(DIR, file)).size;
    const okSize = size < 200 * 1024;
    if (!okSize) bad++;
    console.log(`${okSize ? 'PASS' : 'FAIL'} ${file} size ${(size / 1024).toFixed(1)} KB`);
  }
  fs.writeFileSync(path.join(__dirname, 'verify-report.json'), JSON.stringify(report, null, 1));
  await browser.close();
  console.log(bad ? `\n${bad} FAILED` : '\nALL CLEAN');
  process.exit(bad ? 1 : 0);
})();
