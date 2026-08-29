// Command loadserver serves the same routing table the in-process suite
// attacks, on a real port, for the vegeta CLI to shoot at.
//
// The in-process tests and this binary share app.New(), so a number measured
// through one is a number about the other; what differs is only that the CLI
// can drive it from another machine, for longer, and through vegeta's own
// reporters.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaolaureano/go-router/loadtest/app"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	quiet := flag.Bool("quiet", false, "suppress the route banner")
	catalog := flag.Bool("catalog", false,
		"serve the generated table (three API versions, nested subroutes, mounts) instead of the hand-written one")
	dump := flag.String("dump-targets", "",
		"write a vegeta target file for every route in the table to this path, then exit")
	base := flag.String("base", "http://localhost:8080", "base URL to write into the dumped targets")
	flag.Parse()

	// Dumping targets is how the CLI gets to burst the whole table: the file
	// is generated from the same catalog the server registers, so it cannot
	// drift out of date or miss a route.
	if *dump != "" {
		if err := dumpTargets(*dump, *base); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "wrote %d targets to %s\n", len(app.Catalog()), *dump)
		return
	}

	var handler http.Handler
	if *catalog {
		r, routes := app.NewCatalog()
		handler = r
		if !*quiet {
			fmt.Fprintf(os.Stderr, "go-router load server on %s: generated table, %d routes\n",
				*addr, len(routes))
		}
	} else {
		handler = app.New()
		if !*quiet {
			fmt.Fprintf(os.Stderr, "go-router load server on %s\n\n%s\n", *addr, app.Describe())
		}
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: handler,
		// Deliberately permissive: a load generator opens many connections and
		// holds them, and a timeout firing mid-attack would be measured as the
		// router's latency rather than as the server's policy.
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// A ^C between attacks should report what was served, not just die: the
	// count is the cheapest cross-check that vegeta's request total and the
	// router's handler total agree.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	fmt.Fprintf(os.Stderr, "\nhandled %d requests\n", app.Requests.Load())
}

// dumpTargets renders the catalog as a vegeta target file.
func dumpTargets(path, base string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, route := range app.Catalog() {
		if _, err := fmt.Fprintf(file, "%s %s%s\n", route.Method, base, route.Request); err != nil {
			return err
		}
	}
	return file.Close()
}
