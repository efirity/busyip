# busyip

**Prepaid mobile and residential proxies from real people who are paid per gigabyte.**

[busyip.com](https://busyip.com) · [Pricing](https://busyip.com/pricing) · [Docs](https://busyip.com/docs) · [Start with 1 GB free](https://busyip.com/signup)

---

## Read this before the feature list

busyip is a **small pool of real consumer devices**. Every exit is a person who installed our app,
agreed to share their connection, and is paid for every gigabyte they carry. That is the entire
supply, and it decides what the product can and cannot do:

- **Country coverage is thin and it moves.** Asking for a country with no exit online returns
  `no_capacity` — we refuse rather than quietly route you through somewhere else. Check coverage
  before you promise a country to anyone.
- **Availability is not guaranteed.** These are phones. They go offline when a screen locks, an app
  is backgrounded, or Wi-Fi drops. **Do not point a deadline job at this.**
- **Concurrency is 2 on the free trial and 8 on a paid account.**
- **A sticky session pins a DEVICE, not an IP address.** A carrier can renumber a phone mid-session.
  If your work depends on the address, read it and check it again.

If you need guaranteed uptime across fifty countries, buy a datacentre pool from somebody else. If
you need real consumer connections with a provenance you can explain to a compliance team, and you
can tolerate a fleet that breathes, keep reading.

## Try it before you sign up

**No account needed.** busyip's MCP server answers live capacity, pricing and targeting questions to
anyone:

```bash
curl -s https://busyip.com/mcp \
  -H 'content-type: application/json' \
  -H 'accept: application/json, text/event-stream' \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call",
       "params":{"name":"busyip_capacity","arguments":{}}}'
```

```json
{ "onlineExits": 175,
  "byCountry": { "UZ": 154, "KG": 4, "MD": 4, "DE": 3, "KZ": 2, "RU": 2,
                 "FR": 1, "JP": 1, "SG": 1, "MY": 1, "MX": 1, "TJ": 1 },
  "byType": { "mobile": 112, "residential": 63, "unclassified": 0 },
  "ejectedExits": 1,
  "asOf": "2026-09-09T14:08:36.930Z" }
```

That is a live reading and it will differ when you run it — including the concentration. The fleet
is not evenly spread across those countries and it moves hour to hour, which is why the tool returns
the breakdown rather than a single number you could plan against.

See [MCP.md](MCP.md) for the seven public tools, the six that need OAuth, and how to point an
assistant at them.

## Connect

One credential, HTTP or SOCKS5. **No SDK, no client library** — put the connection string into the
HTTP client you already have.

```
http://USERNAME:PASSWORD@gate.busyip.com:18080
socks5h://USERNAME:PASSWORD@gate.busyip.com:11080
```

Use `socks5h`, not `socks5`. The `h` resolves DNS at the exit instead of on your machine; without it
you leak your own resolver and can defeat the point of the proxy.

Snippets for [curl](examples/curl.sh), [Node](examples/node.mjs), [Python](examples/python.py) and
[Go](examples/main.go).

## Targeting lives in the username

The password never changes. Everything you want to control goes in the username, in this order:

```
USERNAME[-country-{cc}][-type-mobile|residential|any][-session-{id}][-mode-sticky|rotate]
```

| Goal | Username |
|---|---|
| Any exit, new IP each request | `USERNAME` |
| Moldova only | `USERNAME-country-md` |
| Mobile only | `USERNAME-type-mobile` |
| Germany, mobile | `USERNAME-country-de-type-mobile` |
| Hold one device across requests | `USERNAME-session-a1` |
| Force a fresh pick each request | `USERNAME-mode-rotate` |

**A typo matches nothing rather than everything.** `-cuntry-fr` or `-country-zz` returns
`no_capacity`; it never falls back to a different pool. That is deliberate — silently serving you
Germany when you asked for France is worse than refusing.

## What a gigabyte means, and why our number is higher than yours

**1 GB means 1 GiB — 1,073,741,824 bytes — everywhere on this site and in every invoice.**

We count wire bytes in both directions between your client and our gateway, including the TLS
handshake and protocol framing. That reads **above** your application-level counter, and how far
above is set by how much you move per connection rather than by anything we choose:

| Transfer | Overhead | Why |
|---|---|---|
| 120 KB fetch | **+5.6%** | the handshake is amortised over a large body |
| 7 KB fetch | **+46%** | a fresh handshake costs the same however little follows it |

Reusing connections is the single biggest thing that narrows the gap.

**One case where the gap is not small: abandoning a large download.** The phone fetches from the
destination at its own speed, usually faster than a throttled client reads. Cancel part-way and the
phone has already pulled more than reached you — and we count what the phone carried, because that
is the traffic the person was paid for. Ask for a byte range rather than aborting the connection.

**Requests that find no exit cost nothing.** They are refused before any device is touched.

## Pricing

Prepaid. **Bytes never expire**, there is no subscription, and no card is kept on file.

| Pack | Bytes | Price | Per GB |
|---|---|---|---|
| Free trial — 1 GB | 1,073,741,824 | **$0.00** | — |
| 5 GB | 5,368,709,120 | $25.00 | $5.00 |
| 50 GB | 53,687,091,200 | $150.00 | $3.00 |
| 250 GB | 268,435,456,000 | $500.00 | $2.00 |

The free gigabyte is one per person, on signup, **no card**. Buying more adds bytes to the credential
you already have and lifts your connection limit to 8 — **your username and password do not change**,
so the code you got working keeps working.

[Start free →](https://busyip.com/signup)

## Why there is no SDK here

Our sibling product [busysms](https://github.com/efirity/busysms) ships Python and Node SDKs. busyip
deliberately does not, and this is the note that stops someone helpfully adding one.

**A proxy has almost no client surface.** The whole integration is a connection string in an HTTP
client you already have — `requests`, `axios`, `net/http`, curl, your browser automation tool, your
scraper. A busyip SDK would wrap one string and then need maintaining, versioning and releasing.

**The part that could use a helper is the targeting grammar**, and that is a string format. Wrapping
it in a package would put a version number between you and a contract that is enforced by the
gateway, not by us — a misspelled marker is refused at the gate whatever your client library thinks.
You are better off with the table above and a `f"{user}-country-{cc}"` in your own code, where you
can see what is being sent.

If you want programmatic access to *capacity, pricing, or your account*, that is what
[the MCP server](MCP.md) is for, and it needs no library either.

## Where the exits come from

Every exit is a person running our app who chose to share their connection and is paid per gigabyte.
No SDK bundled into someone else's software, no traffic borrowed from users who did not read what
they agreed to.

That is also why the pool is small, and we would rather say so here than have you discover it in a
benchmark.

## Acceptable use

Proxies are dual-use and we are not neutral about it. Credential stuffing, fraud, spam, scraping
behind a login you do not own, and anything targeting a private individual are out. Port access is
**80 and 443 only** — no SMTP, no arbitrary ports — and the gateway refuses connections to private
and loopback address ranges, so the fleet cannot be pointed at an earner's own home network.

Full terms: [busyip.com/legal](https://busyip.com/legal/terms) · Seller: **EFIRITY PTE. LTD.**
(Singapore, UEN 202402844G)

## Support

[support@busyip.com](mailto:support@busyip.com). If your byte figures look wrong even allowing for
the overhead above, tell us — we would rather check and be wrong than have you assume.
