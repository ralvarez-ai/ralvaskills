// node smoke.js direction-1.html [...]
// Every control clicked at 1440 and 390; fail on pageerror; axe re-run in the end state.
const { launch, open, axe, overflow } = require('./lib');

(async () => {
  const files = process.argv.slice(2);
  const browser = await launch();
  let bad = 0;
  const fail = (m) => { bad++; console.log('  FAIL', m); };
  for (const file of files) for (const width of [1440, 390]) {
    const { ctx, page, errors } = await open(browser, file, { width, theme: 'light', reduced: false });
    await ctx.grantPermissions(['clipboard-read', 'clipboard-write']).catch(() => {});
    const log = [];
    // skip link is off-screen until focused: drive it by keyboard
    await page.keyboard.press('Tab');
    const first = await page.evaluate(() => document.activeElement.className);
    if (first !== 'skip') fail('first Tab stop is ' + first);
    await page.keyboard.press('Enter');
    if (!(await page.evaluate(() => location.hash === '#main'))) fail('skip link did not move to #main');
    await page.evaluate(() => { history.replaceState(null, '', location.pathname); document.documentElement.style.scrollBehavior = 'auto'; scrollTo(0, 0); });
    // films: pause, resume, let it settle, replay
    const films = await page.$$('[data-film]');
    for (const [i, f] of films.entries()) {
      await f.scrollIntoViewIfNeeded();
      await page.waitForFunction((el) => el.kFilm && el.kFilm.state() === 'playing', f, { timeout: 4000 }).catch(() => fail(`film ${i} did not autoplay`));
      const ctl = await f.$('.k-film-ctl');
      await ctl.click();
      if ((await f.evaluate((el) => el.kFilm.state())) !== 'paused') fail(`film ${i} pause`);
      const label = await ctl.innerText();
      if (!/Play/.test(label)) fail(`film ${i} label after pause: ${label}`);
      await ctl.click();
      await page.waitForFunction((el) => el.kFilm.state() === 'done', f, { timeout: 12000 }).catch(() => fail(`film ${i} never settled`));
      await ctl.click();
      if ((await f.evaluate((el) => el.kFilm.state())) !== 'playing') fail(`film ${i} replay`);
      await page.waitForFunction((el) => el.kFilm.state() === 'done', f, { timeout: 12000 }).catch(() => fail(`film ${i} replay never settled`));
      log.push(`film${i} ok`);
    }
    // theme toggle (hidden labels at 390 still click)
    const tb = await page.$('.k-theme');
    for (const want of ['dark', 'light', 'dark']) {
      await tb.click();
      await page.waitForTimeout(650);
      const t = await page.evaluate(() => [document.documentElement.dataset.theme, document.querySelector('.k-theme').getAttribute('aria-pressed'), localStorage.getItem('rsk-theme')]);
      if (t[0] !== want || t[1] !== String(want === 'dark') || t[2] !== want) fail(`theme toggle -> ${t}`);
    }
    log.push('theme ok');
    // copy buttons
    for (const b of await page.$$('[data-copy]')) {
      await b.scrollIntoViewIfNeeded();
      await b.click();
      await page.waitForTimeout(150);
      const txt = await b.innerText();
      if (!/Copied|Copy failed/.test(txt)) fail('copy label ' + txt);
      const clip = await page.evaluate(() => navigator.clipboard.readText().catch(() => ''));
      if (clip && !clip.startsWith('go install')) fail('clipboard content ' + clip.slice(0, 30));
    }
    log.push('copy ok');
    // catalogue filter
    if (await page.$('#b-chips')) {
      const count = () => page.$eval('#b-count', (e) => +e.textContent);
      for (const [bundle, want] of [['go-grpc', 5], ['llm-app', 4], ['docs', 0], ['global', 15], ['', 49]]) {
        const c = await page.$(`.b-chip[data-bundle="${bundle}"]`);
        await c.scrollIntoViewIfNeeded();
        await c.click();
        const n = await count();
        if (n !== want) fail(`bundle ${bundle || 'all'} count ${n} != ${want}`);
      }
      for (const c of await page.$$('.b-chip')) await c.click();
      const inc = await page.$('#b-inc');
      if ((await inc.getAttribute('aria-pressed')) !== 'true') await inc.click();
      await page.click('.b-chip[data-bundle=""]');
      if ((await count()) !== 58) fail('include personal count ' + (await count()));
      await page.click('.b-chip[data-bundle="code-review"]');
      log.push('filter ok');
    }
    // details + in-page links
    for (const s of await page.$$('details > summary')) { await s.scrollIntoViewIfNeeded(); await s.click(); }
    const hrefs = await page.$$eval('a[href^="#"]', (as) => as.filter((a) => a.offsetParent).map((a) => a.getAttribute('href')));
    for (const h of [...new Set(hrefs)].filter((x) => x !== '#main')) {
      if (!(await page.$(h === '#' ? 'body' : h))) fail('missing target ' + h);
      await page.click(`a[href="${h}"]:visible >> nth=0`).catch((e) => fail('click ' + h + ' ' + e.message.split('\n')[0]));
    }
    log.push(`${new Set(hrefs).size} links ok`);
    // any other button not yet covered
    const n = await page.$$eval('button', (bs) => bs.filter((b) => !b.hidden && !b.matches('.k-theme,.k-film-ctl,[data-copy],.b-chip')).length);
    if (n) fail(`${n} uncovered buttons`);
    // end state: dark, filtered, details open -> axe again
    await page.evaluate(() => { document.documentElement.style.scrollBehavior = 'auto'; scrollTo(0, 0); });
    await page.waitForTimeout(400);
    const v = await axe(page);
    if (v.length) fail('axe end state: ' + v.map((x) => `${x.id}(${x.n}) ${x.targets[0]}`).join(' | '));
    const o = await overflow(page);
    if (o.overflow) fail('overflow end state ' + o.culprits.join(','));
    if (errors.length) fail('errors: ' + errors.join(' | '));
    console.log(`${file} ${width}: ${log.join(', ')}`);
    await ctx.close();
  }
  await browser.close();
  console.log(bad ? `\n${bad} FAILED` : '\nSMOKE CLEAN');
  process.exit(bad ? 1 : 0);
})();
