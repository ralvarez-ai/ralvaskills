// The landing is static assets; this Worker runs first only for /registry/*, where it serves
// the registry index under the landing host so `RSK_REGISTRY_URL=https://<landing host>/registry`
// works. The registry's own host is configuration (REGISTRY_HOST, bound as REGISTRY_ORIGIN by
// scripts/deploy.mjs) and stays the one publisher; nothing is copied or cached here beyond
// what the edge does. Only GET/HEAD of index.json is proxied: the index carries absolute
// archive URLs, so tarballs never pass through this Worker.
const PREFIX = '/registry/'
const ALLOWED = new Set(['index.json'])

export default {
	async fetch(request, env) {
		const url = new URL(request.url)
		if (!url.pathname.startsWith(PREFIX)) {
			return env.ASSETS.fetch(request)
		}
		if (request.method !== 'GET' && request.method !== 'HEAD') {
			return new Response('method not allowed', { status: 405, headers: { Allow: 'GET, HEAD' } })
		}
		const file = url.pathname.slice(PREFIX.length)
		if (!ALLOWED.has(file)) {
			return new Response('not found', { status: 404 })
		}
		if (!env.REGISTRY_ORIGIN) {
			return new Response('REGISTRY_ORIGIN is not set', { status: 500 })
		}
		const upstream = await fetch(`${env.REGISTRY_ORIGIN.replace(/\/+$/, '')}/${file}`, {
			method: request.method,
			headers: { Accept: 'application/json' },
			cf: { cacheTtl: 300, cacheEverything: true },
		})
		if (!upstream.ok) {
			return new Response('registry unavailable', { status: 502 })
		}
		return new Response(upstream.body, {
			status: 200,
			headers: {
				'Content-Type': 'application/json; charset=utf-8',
				'Cache-Control': 'public, max-age=300',
			},
		})
	},
}
