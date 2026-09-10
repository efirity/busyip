// busyip from Node. No dependency on busyip — only an HTTP agent that speaks CONNECT.
//
//   npm i undici
//   BUSYIP_USER=... BUSYIP_PASS=... node node.mjs
//
// Credentials from the environment, never from argv: a password on a command line is visible in
// `ps` to every user on the box. Get a free gigabyte at https://busyip.com/signup

import { ProxyAgent, request } from 'undici'

const USER = process.env.BUSYIP_USER
const PASS = process.env.BUSYIP_PASS
if (!USER || !PASS) throw new Error('set BUSYIP_USER and BUSYIP_PASS')

const GATE = 'http://gate.busyip.com:18080'

/**
 * Build a targeted username.
 *
 * The marker ORDER is fixed by the gateway's grammar, which is why this composes rather than
 * concatenating in call order. A misspelled marker NAME is refused at the gate with no_capacity —
 * it is not ignored — so a typo here fails loudly rather than quietly widening your pool.
 */
function target (user, { country, type, session, mode } = {}) {
  let u = user
  if (country) u += `-country-${country}`
  if (type) u += `-type-${type}`
  if (session) u += `-session-${session}`
  if (mode) u += `-mode-${mode}`
  return u
}

const via = (opts) =>
  new ProxyAgent({ uri: GATE, token: `Basic ${Buffer.from(`${target(USER, opts)}:${PASS}`).toString('base64')}` })

const exitIp = async (opts) => {
  const res = await request('https://icanhazip.com', { dispatcher: via(opts) })
  return (await res.body.text()).trim()
}

console.log('any exit      :', await exitIp())
console.log('germany       :', await exitIp({ country: 'de' }))
console.log('mobile only   :', await exitIp({ type: 'mobile' }))

// A session pins a DEVICE, not an address: the same phone answers, and a carrier may still
// renumber it between requests. Read the address if your work depends on it.
const sid = `job-${process.pid}`
console.log('session, 1st  :', await exitIp({ session: sid }))
console.log('session, 2nd  :', await exitIp({ session: sid }))

// A country with nobody online refuses rather than substituting somewhere else.
try { console.log('country-zz    :', await exitIp({ country: 'zz' })) }
catch (err) { console.log('country-zz    : refused —', err.message) }
