package httpserver

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jaeger-tenant/pkg/helm"
	"jaeger-tenant/pkg/jaeger"
	"jaeger-tenant/pkg/structures"
	"jaeger-tenant/utils"
	"log"
	"net/http"
	"time"
)

// StartHTTPServer start the http server
func StartHTTPServer() {
	helm.DownloadHelmChart()
	muxHTTP := mux.NewRouter()
	muxHTTP.HandleFunc("/tenant", TenantPut).Methods("PUT")
	muxHTTP.HandleFunc("/tenant/{customer}", TenantGet).Methods("GET")
	muxHTTP.HandleFunc("/tenant/{customer}", TenantDelete).Methods("DELETE")
	muxHTTP.HandleFunc("/tenant", TenantGetList).Methods("GET")
	muxHTTP.HandleFunc("/performance", PerformancePut).Methods("PUT")
	fmt.Println("Starting server on Port 8080")

	srv := &http.Server{
		Handler: muxHTTP,
		Addr:    "0.0.0.0:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Httpserver: ListenAndServe() error: %s", err)
	}
}

// TenantPut handles the put request of Tenant
func TenantPut(w http.ResponseWriter, r *http.Request) {
	var payload structures.TenantPayload
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	var allTenantDetails []structures.Tenant

	for _, customer := range payload.Customers {

		allTenantDetails = append(allTenantDetails, jaeger.CreateTenant(customer, &payload))
	}
	result, _ := json.Marshal(allTenantDetails)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(result)
}

// TenantGet returns the JSON containing details of Tenant
func TenantGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer := vars["customer"]

	b := jaeger.GetTenant(customer)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(b)
}

// TenantDelete deletes the tenant, returning a message
func TenantDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer := vars["customer"]

	b := jaeger.DeleteTenant(customer)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(b)
}

// TenantGetList returns a JSON containing list of all Tenants
func TenantGetList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

}

// PerformancePut for creating performance test on Tenants
func PerformancePut(w http.ResponseWriter, r *http.Request) {
	var payload structures.PerformancePayload
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&payload)
	if err != nil {
		panic(err)
	}

	var b []byte
	err = jaeger.CreatePerformanceTest(&payload)
	if err != nil {
		b, _ = json.Marshal(utils.CreateErrorResponse(err))
	}

	b, _ = json.Marshal(utils.CreateMessageResponse("success"))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(b)
}
