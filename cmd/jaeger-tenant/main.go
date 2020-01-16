package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func main() {
	http.HandleFunc("/", handlerHelloWorld)
	http.ListenAndServe(":8080", nil)
}

func handlerHelloWorld(w http.ResponseWriter, r *http.Request) {

	helmArg := "helm template --set storage.keyspace=arpit,storage.cassandra.host=chart-1577686576-cassandra,storage.cassandra.dc_name=dc1,storage.cassandra.password=password,storage.kafka.enabled=true,storage.kafka.broker=my-kafka:9092 my-jaeger /Users/I353342/workspace/helm/jaeger/"

	out, err := exec.Command("/bin/bash", "-c", helmArg).Output()
	if err != nil {
		fmt.Print(err)
	}
	fmt.Fprintf(w, string(out[:]))
}
