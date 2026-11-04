// node shots.js direction-1.html [...]  -> ../shots/<stem>-<width>-<theme>[-reduced|-full].png
// Settled film (motion) first viewport + full page, and reduced motion first viewport.
const path = require('path');
const { DIR, launch, open, settle } = require('./lib');

(async () => {
  const files = process.argv.slice(2);
  const only = process.env.ONLY; // e.g. "fold" to skip full pages
  const browser = await launch();
  for (const file of files) {
    const stem = path.basename(file, '.html');
    for (const width of [1440, 390]) for (const theme of ['light', 'dark']) {
      for (const reduced of [false, true]) {
        const { ctx, page, errors } = await open(browser, file, { width, theme, reduced });
        await settle(page, { reduced });
        const base = path.join(DIR, 'shots', `${stem}-${width}-${theme}${reduced ? '-reduced' : ''}`);
        await page.screenshot({ path: base + '.png' });
        if (!reduced && only !== 'fold') await page.screenshot({ path: base + '-full.png', fullPage: true });
        if (errors.length) console.log('ERR', file, errors);
        await ctx.close();
      }
    }
    console.log('shots', file);
  }
  await browser.close();
})();
