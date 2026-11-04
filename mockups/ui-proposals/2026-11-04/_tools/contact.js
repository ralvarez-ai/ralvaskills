// node contact.js stem1 stem2 ... -> ../shots/contact-sheet.png (first viewports, 1440 + 390, light + dark)
const path = require('path');
const { DIR, launch } = require('./lib');
(async () => {
  const stems = process.argv.slice(2);
  const out = process.env.OUT || 'contact-sheet.png';
  const S = (f) => 'data:image/png;base64,' + require('fs').readFileSync(path.join(DIR, 'shots', f)).toString('base64');
  const rows = stems.map((s) => `<tr><th>${s}</th>${['light', 'dark'].map((t) => `<td><img src="${S(`${s}-1440-${t}.png`)}" width="720"></td><td><img src="${S(`${s}-390-${t}.png`)}" width="195"></td>`).join('')}</tr>`).join('');
  const html = `<html><body style="margin:0;background:#888;font:16px sans-serif"><table cellspacing="12"><tr><th></th><th>1440 light</th><th>390 light</th><th>1440 dark</th><th>390 dark</th></tr>${rows}</table></body></html>`;
  const b = await launch();
  const p = await b.newPage({ viewport: { width: 1900, height: 600 } });
  await p.setContent(html, { waitUntil: 'load' });
  await p.screenshot({ path: path.join(DIR, 'shots', out), fullPage: true });
  await b.close();
})();
