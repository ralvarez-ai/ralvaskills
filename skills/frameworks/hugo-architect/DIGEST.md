# hugo-architect (digest)

Architecture standards for Hugo 0.161 Extended sites and themes: layout, front matter, templates, modules, page bundles, Hugo Pipes, i18n, render hooks, deploy.

## Always-on rules

- Use Hugo Extended 0.161; Sass and image processing need it.
- `content/` mirrors URLs; `assets/` holds processed files (Hugo Pipes); `static/` holds pass-through files; `data/` holds structured data, never content.
- One front-matter format per repo: TOML is the default. Required on every page: `title`, `date`, `draft`.
- Taxonomies are arrays. Set `slug` only when the title is unstable. Use `weight` only for manual ordering. Custom keys go under `[params]`.
- Shared themes ship as Hugo Modules, not submodules. Vendor with `hugo mod vendor` and commit `_vendor/`.
- Use page bundles (`index.md` plus assets) as soon as a page has images.
- Pipeline: toCSS or js.Build, then postCSS if used, minify, fingerprint. Fingerprinting is mandatory in production.
- Use `i18n` keys for UI strings. Decide the default-language URL prefix early.
- Prefer render hooks over shortcodes; shortcodes are for structured content, never styling.
- No JS by default. Build with `--minify`.
- Set `baseURL` to the deploy domain; CI injects it per environment.
- Static CDN hosting by default (Cloudflare Pages, Netlify, GitHub Pages, S3 + CloudFront).

## Hard constraints

- Never inject styling through shortcodes.
- Never hardcode UI strings in templates.
- Never mix TOML and YAML front matter in one repo.
- Never rely on the global scope inside a partial; pass an explicit context.
- Enable `unsafe = true` only when content authors are trusted.

Invoke the full skill when: scaffolding, reviewing or auditing a Hugo site or theme, or choosing modules, bundles, pipes, i18n or deploy setup.

Full guide: invoke the hugo-architect skill.
