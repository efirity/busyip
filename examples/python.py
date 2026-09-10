"""busyip from Python. Only `requests` — there is no busyip package to install.

    pip install requests
    BUSYIP_USER=... BUSYIP_PASS=... python python.py

Credentials from the environment, never from argv: a password on a command line is visible in `ps`
to every user on the box. Get a free gigabyte at https://busyip.com/signup
"""
import os
import requests

USER = os.environ["BUSYIP_USER"]
PASS = os.environ["BUSYIP_PASS"]
HOST = "gate.busyip.com"


def target(user, country=None, type=None, session=None, mode=None):
    """Compose a targeted username.

    The marker ORDER is fixed by the gateway's grammar. A misspelled marker NAME is refused with
    no_capacity rather than ignored, so a typo fails loudly instead of quietly widening your pool.
    """
    u = user
    if country:
        u += f"-country-{country}"
    if type:
        u += f"-type-{type}"
    if session:
        u += f"-session-{session}"
    if mode:
        u += f"-mode-{mode}"
    return u


def proxies(**opts):
    uri = f"http://{target(USER, **opts)}:{PASS}@{HOST}:18080"
    return {"http": uri, "https": uri}


def socks(**opts):
    # socks5h, not socks5 — the `h` resolves DNS at the exit rather than on this machine.
    # Needs: pip install "requests[socks]"
    uri = f"socks5h://{target(USER, **opts)}:{PASS}@{HOST}:11080"
    return {"http": uri, "https": uri}


def exit_ip(**opts):
    return requests.get("https://icanhazip.com", proxies=proxies(**opts), timeout=60).text.strip()


print("any exit     :", exit_ip())
print("germany      :", exit_ip(country="de"))
print("mobile only  :", exit_ip(type="mobile"))

# A session pins a DEVICE, not an address. The same phone answers; its carrier may still renumber it.
sid = f"job-{os.getpid()}"
print("session, 1st :", exit_ip(session=sid))
print("session, 2nd :", exit_ip(session=sid))

# Reuse the connection and the metering overhead drops sharply — see the README's byte table.
with requests.Session() as s:
    s.proxies.update(proxies(country="de"))
    for _ in range(3):
        s.get("https://icanhazip.com", timeout=60)

# A country with nobody online refuses rather than substituting somewhere else.
try:
    print("country-zz   :", exit_ip(country="zz"))
except requests.RequestException as err:
    print("country-zz   : refused —", err)
