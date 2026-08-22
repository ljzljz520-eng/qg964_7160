package main

import (
	"coldrisk.local/console/internal/domain"
	"coldrisk.local/console/internal/httpapi"
	"coldrisk.local/console/internal/policy"
	"coldrisk.local/console/internal/service"
	"coldrisk.local/console/internal/store"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	path := flag.String("db", "risk-console.db", "embedded bbolt database")
	addr := flag.String("addr", ":8080", "http listen address")
	flag.Parse()
	st, err := store.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer st.Close()
	pol := defaultPolicy()
	svc := service.New(st, pol)
	server := httpapi.New(svc)
	fmt.Println("cold storage risk console listening on", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
func defaultPolicy() *policy.Policy {
	p := policy.New()
	p.Set(policyPermission("demo-supervisor", []string{"store-a", "store-b"}))
	return p
}
func policyPermission(principal string, stores []string) domain.Permission {
	actions := map[string]bool{"view_todo": true, "create_record": true, "review_record": true, "publish_record": true, "archive_record": true, "import_records": true}
	scope := map[string]string{}
	for _, storeID := range stores {
		scope[storeID] = "regional"
	}
	return domain.Permission{Principal: principal, Roles: []string{"regional_supervisor"}, Stores: scope, Actions: actions}
}
