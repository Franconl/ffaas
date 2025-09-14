package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Franconl/ffaas/internal/httpapi"
	"github.com/Franconl/ffaas/internal/repo/memory"
)

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	store := memory.New()
	handler := httpapi.NewRouter(store)

	log.Printf("faas-server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
