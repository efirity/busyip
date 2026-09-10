# busyip over MCP

busyip serves a [Model Context Protocol](https://modelcontextprotocol.io) endpoint at
**`https://busyip.com/mcp`**. Seven tools answer without any account at all; six more act on a
signed-in customer's account and require OAuth.

This exists because a proxy is a thing an agent buys and configures, not just a thing a human clicks.
An assistant can check whether a country is online, quote a price, and build a correctly-targeted
username without a key — and only then send its user to sign up.

## Add it to a client

**Claude Code**

```bash
claude mcp add --transport http busyip https://busyip.com/mcp
```

**Anything reading `mcpServers` config**

```json
{ "mcpServers": { "busyip": { "type": "http", "url": "https://busyip.com/mcp" } } }
```

The public tools work immediately. The account tools trigger an OAuth flow the first time one is
called — busyip runs its own OAuth 2.1 server with PKCE and dynamic client registration, so there is
no key to paste and nothing to store.

## `tools/list` shows all thirteen — the gate is on calling, not discovery

Run `tools/list` with no credentials and you get **all 13 tools**, including the six that need
OAuth. That is deliberate: a tool an assistant cannot see is a tool it cannot offer, and hiding the
account tools from an anonymous probe would mean no client could ever register them.

Calling a gated tool without a token is **not** a protocol error. It returns a normal result
carrying `error` and `hint`:

```json
{ "error": "unauthorized: no token",
  "hint": "\"busyip_balance\" needs a busyip MCP token … The public tools work without one." }
```

So a client exploring the server never hits a hard failure — it gets a readable sentence telling it
what to do next. Advertising is not exposing: `tools/call` still refuses without a verified customer.

## Public tools — no account

| Tool | Answers |
|---|---|
| `busyip_overview` | what the product is, the endpoint, and its limits |
| `busyip_capacity` | **live** exits online, by country and by network type |
| `busyip_pricing` | packs, per-GB rates, the free trial |
| `busyip_targeting` | the username grammar, with worked examples |
| `busyip_metering` | what a byte is, and why our count reads above yours |
| `busyip_legal` | seller entity, terms, acceptable use |
| `busyip_brand` | name, icons, colours — for anyone integrating busyip into their own product |

## Account tools — OAuth

| Tool | Does |
|---|---|
| `busyip_balance` | bytes remaining, used, granted |
| `busyip_usage` | per-day usage |
| `busyip_credential` | the proxy username and connection strings |
| `busyip_orders` | purchase history and invoices |
| `busyip_account` | account status |
| `busyip_rotate_secret` | mints a new password — **requires `confirm: true`** |

`busyip_rotate_secret` is deliberately awkward. Rotating replaces the password your running code is
using, so it cannot happen as a side effect of an agent exploring the tool list.

## Try a tool with no client

```bash
curl -s https://busyip.com/mcp \
  -H 'content-type: application/json' \
  -H 'accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'
```

Then call one:

```bash
curl -s https://busyip.com/mcp \
  -H 'content-type: application/json' \
  -H 'accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call",
       "params":{"name":"busyip_targeting","arguments":{}}}'
```

## What `busyip_capacity` does and does not tell you

It reports exits that are **online, exit-enabled and not currently ejected for failing**. That is
supply availability, not a promise that your next request succeeds — a device can go offline between
your read and your request, and the fleet turns over roughly a third of its online set within an
hour.

Two fields are there so you can tell the difference between "the fleet changed" and "something is
wrong": `ejectedExits` counts exits removed for failing health, and `unclassified` counts exits the
router cannot place in either the mobile or residential pool. **`byCountry` and `byType` each sum to
`onlineExits`** — if they ever do not, the breakdown is lying and we would like to hear about it.

Read it before a run and after it. If `ejectedExits` moved, the fleet changed underneath you and a
failed country is not evidence of a routing bug.

---

Questions: [support@busyip.com](mailto:support@busyip.com) · [Sign up](https://busyip.com/signup)
