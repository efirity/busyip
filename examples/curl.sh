#!/usr/bin/env bash
# busyip — the whole integration. No SDK, no library.
#
# Get a username and password free at https://busyip.com/signup (1 GB, no card).
# Never put the password on a command line on a shared machine: it is visible in `ps` to every user.
set -euo pipefail

URI="${BUSYIP_URI:?export BUSYIP_URI=http://USERNAME:PASSWORD@gate.busyip.com:18080}"
USER_PART="${URI#*://}"; USER_PART="${USER_PART%%:*}"
REST="${URI#*://*:}"

# 1. Plainest possible request.
curl -s -x "$URI" https://icanhazip.com

# 2. Target a country. The marker goes in the USERNAME, not in a header or a flag.
curl -s -x "http://${USER_PART}-country-de:${REST}" https://icanhazip.com

# 3. Hold one device across several requests: same session id, same phone.
#    A session pins a DEVICE, not an address — a carrier can renumber it mid-session.
SID="job-$$"
for i in 1 2 3; do
  curl -s -x "http://${USER_PART}-session-${SID}:${REST}" https://icanhazip.com
done

# 4. SOCKS5 — note socks5h. The `h` resolves DNS at the exit rather than on your machine.
curl -s -x "socks5h://${USER_PART}:${REST/18080/11080}" https://icanhazip.com

# 5. Ask for a country with nobody online and you get a refusal, not a silent substitution.
curl -s -o /dev/null -w 'no-capacity case: HTTP %{http_code} in %{time_total}s\n' \
  -x "http://${USER_PART}-country-zz:${REST}" https://icanhazip.com || true
