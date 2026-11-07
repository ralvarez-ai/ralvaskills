import assert from 'node:assert/strict'
import { afterEach, test } from 'node:test'
import worker from './index.js'

const assets = { fetch: async () => new Response('asset') }
const env = { ASSETS: assets, REGISTRY_ORIGIN: 'https://registry.example' }
const realFetch = globalThis.fetch
afterEach(() => {
	globalThis.fetch = realFetch
})
const call = (path, init, vars = env) => worker.fetch(new Request(`https://landing.example${path}`, init), vars)

test('other paths go to the static assets', async () => {
	assert.equal(await (await call('/')).text(), 'asset')
})

test('/registry/index.json is proxied from the configured origin', async () => {
	let asked
	globalThis.fetch = async (u) => {
		asked = String(u)
		return new Response('{"skills":[]}', { status: 200 })
	}
	const res = await call('/registry/index.json')
	assert.equal(asked, 'https://registry.example/index.json')
	assert.equal(res.status, 200)
	assert.equal(res.headers.get('content-type'), 'application/json; charset=utf-8')
})

test('only index.json is served under /registry/', async () => {
	assert.equal((await call('/registry/other.tgz')).status, 404)
	assert.equal((await call('/registry/')).status, 404)
})

test('methods other than GET and HEAD are refused', async () => {
	assert.equal((await call('/registry/index.json', { method: 'POST' })).status, 405)
})

test('without REGISTRY_ORIGIN it fails closed', async () => {
	assert.equal((await call('/registry/index.json', undefined, { ASSETS: assets })).status, 500)
})

test('an upstream failure is a 502', async () => {
	globalThis.fetch = async () => new Response('x', { status: 500 })
	assert.equal((await call('/registry/index.json')).status, 502)
})
