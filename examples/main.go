// busyip from Go. Standard library only — there is no busyip module to import.
//
//	BUSYIP_USER=... BUSYIP_PASS=... go run main.go
//
// Credentials from the environment, never from argv: a password on a command line is visible in
// `ps` to every user on the box. Get a free gigabyte at https://busyip.com/signup
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const host = "gate.busyip.com:18080"

type opts struct {
	country string
	typ     string
	session string
	mode    string
}

// target composes a targeted username.
//
// The marker ORDER is fixed by the gateway's grammar, which is why this builds in sequence rather
// than in call order. A misspelled marker NAME is refused at the gate with no_capacity — it is not
// ignored — so a typo fails loudly rather than quietly widening the pool you are drawing from.
func target(user string, o opts) string {
	var b strings.Builder
	b.WriteString(user)
	if o.country != "" {
		b.WriteString("-country-" + o.country)
	}
	if o.typ != "" {
		b.WriteString("-type-" + o.typ)
	}
	if o.session != "" {
		b.WriteString("-session-" + o.session)
	}
	if o.mode != "" {
		b.WriteString("-mode-" + o.mode)
	}
	return b.String()
}

func client(user, pass string, o opts) *http.Client {
	proxy := &url.URL{
		Scheme: "http",
		User:   url.UserPassword(target(user, o), pass),
		Host:   host,
	}
	return &http.Client{
		Timeout:   60 * time.Second,
		Transport: &http.Transport{Proxy: http.ProxyURL(proxy)},
	}
}

func exitIP(user, pass string, o opts) (string, error) {
	res, err := client(user, pass, o).Get("https://icanhazip.com")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(body)), nil
}

func main() {
	user, pass := os.Getenv("BUSYIP_USER"), os.Getenv("BUSYIP_PASS")
	if user == "" || pass == "" {
		fmt.Fprintln(os.Stderr, "set BUSYIP_USER and BUSYIP_PASS")
		os.Exit(1)
	}

	show := func(label string, o opts) {
		ip, err := exitIP(user, pass, o)
		if err != nil {
			// A country with nobody online is refused rather than substituted. That is the
			// designed behaviour, not a fault: handle it as "no supply", not as an outage.
			fmt.Printf("%-14s refused — %v\n", label, err)
			return
		}
		fmt.Printf("%-14s %s\n", label, ip)
	}

	show("any exit", opts{})
	show("germany", opts{country: "de"})
	show("mobile only", opts{typ: "mobile"})

	// A session pins a DEVICE, not an address. The same phone answers; its carrier may still
	// renumber it between requests, so read the address when your work depends on it.
	sid := fmt.Sprintf("job-%d", os.Getpid())
	show("session, 1st", opts{session: sid})
	show("session, 2nd", opts{session: sid})

	show("country-zz", opts{country: "zz"})
}
