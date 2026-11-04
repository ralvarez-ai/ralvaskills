// node sec.js file width theme sel1 sel2 ... -> ../shots/sec-<stem>-<width>-<theme>-<sel>.png (settled)
const path = require('path');
const { DIR, launch, open, settle } = require('./lib');
(async () => {
  const [file, width, theme, ...sels] = process.argv.slice(2);
  const b = await launch();
  const { ctx, page } = await open(b, file, { width: +width, theme, reduced: false });
  await settle(page, { reduced: false });
  for (const s of sels) {
    const el = await page.$(s);
    await el.screenshot({ path: path.join(DIR, 'shots', `sec-${path.basename(file, '.html')}-${width}-${theme}-${s.replace(/[^a-z0-9-]/gi, '')}.png`) });
  }
  await b.close();
})();
