/* Mockup preview bar. Declarative: set window.PV before loading, or put the JSON in <script type="application/json" id="pv-config">.
   PV = { axes: [ { key:'palette', label:'Palette', values:[{v:'a',label:'A teal'},{v:'b',label:'B ink'}],
                    when:{key:'palette',is:'b'} (optional: show only while another axis has this value) } ] }
   A value may carry `sw` (swatch colour) and `swDark` (swatch in data-theme=dark).
   The first value of each axis is the default. State -> <html data-KEY=value>, localStorage 'pv-KEY', and ?KEY=value
   (the query wins, is not persisted, and is how a screenshot script forces an exact variant). Run in <head> to avoid a flash. */
(function () {
  var cfg = window.PV || JSON.parse((document.getElementById('pv-config') || {}).textContent || '{"axes":[]}');
  var R = document.documentElement, axes = cfg.axes || [];
  function ok(a, v) { return a.values.some(function (x) { return x.v === v; }); }
  function stored(k) { try { return localStorage.getItem('pv-' + k); } catch (e) { return null; } }
  function pick(a) {
    var q = new URLSearchParams(location.search).get(a.key), s = stored(a.key);
    return ok(a, q) ? q : ok(a, s) ? s : a.values[0].v;
  }
  function set(a, v, keep) {
    R.setAttribute('data-' + a.key, v);
    if (!keep) { try { localStorage.setItem('pv-' + a.key, v); } catch (e) {} }
    [].forEach.call(document.querySelectorAll('.pv [data-set="' + a.key + '"]'), function (b) {
      b.setAttribute('aria-pressed', b.getAttribute('data-v') === v ? 'true' : 'false');
    });
    gate();
  }
  function gate() {
    axes.forEach(function (a) {
      var f = document.querySelector('.pv [data-axis="' + a.key + '"]');
      if (f && a.when) f.hidden = R.getAttribute('data-' + a.when.key) !== a.when.is;
    });
  }
  function build() {
    var bar = document.createElement('div');
    bar.className = 'pv'; bar.setAttribute('role', 'region'); bar.setAttribute('aria-label', 'Mockup preview controls');
    axes.forEach(function (a) {
      var f = document.createElement('fieldset'), l = document.createElement('legend');
      f.setAttribute('data-axis', a.key); l.textContent = a.label; f.appendChild(l);
      a.values.forEach(function (x) {
        var b = document.createElement('button');
        b.type = 'button'; b.setAttribute('data-set', a.key); b.setAttribute('data-v', x.v); b.setAttribute('aria-pressed', 'false');
        if (x.sw) { var i = document.createElement('i'); i.className = 'sw'; i.setAttribute('aria-hidden', 'true'); i.style.setProperty('--sw', x.sw); if (x.swDark) i.style.setProperty('--sw-dark', x.swDark); b.appendChild(i); }
        b.appendChild(document.createTextNode(x.label)); f.appendChild(b);
      });
      bar.appendChild(f);
    });
    document.body.appendChild(bar);
  }
  axes.forEach(function (a) { R.setAttribute('data-' + a.key, pick(a)); });
  function mount() {
    build();
    axes.forEach(function (a) { set(a, pick(a), true); });
    document.addEventListener('click', function (e) {
      var b = e.target.closest && e.target.closest('.pv [data-set]');
      if (!b) return;
      var a = axes.filter(function (x) { return x.key === b.getAttribute('data-set'); })[0];
      if (a) set(a, b.getAttribute('data-v'));
    });
  }
  if (document.body) mount(); else document.addEventListener('DOMContentLoaded', mount);
})();
