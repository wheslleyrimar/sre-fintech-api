package main

import (
	"fmt"
	stdhttp "net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"sre/internal/http"
	httpclient "sre/internal/httpClient"
	"sre/internal/integrations"
	"sre/internal/usecases"
)

func main() {
	backendURL := "http://localhost:8080"
	if u := os.Getenv("BACKEND_URL"); u != "" {
		backendURL = u
	}
	myselfURL := "http://localhost:8080"
	if u := os.Getenv("SRE_BASE_URL"); u != "" {
		myselfURL = u
	}
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	factory := httpclient.NewEndpointFactory(backendURL)
	searchEngine := integrations.NewSearchEngine(factory)
	accountsAPI := integrations.NewAccountsApi(factory)
	adjustmentFlow := integrations.NewAdjustmentFlowProcessor(factory)

	searchSvc := usecases.NewSearchService(searchEngine)
	accountSvc := usecases.NewAccountService(accountsAPI, accountsAPI, adjustmentFlow, myselfURL)
	reportSvc := usecases.NewReportService(searchSvc)

	r := chi.NewRouter()
	r.Route("/v1", func(r chi.Router) {
		http.NewAccountController(accountSvc).Routes(r)
		http.NewReportController(reportSvc).Routes(r)
		http.NewSearchController(searchSvc).Routes(r)
	})

	fmt.Printf("SRE API listening on http://localhost%s (backend: %s)\n", addr, backendURL)
	if err := stdhttp.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}
