package utils

import (
	"github.com/johnaoss/htpasswd/apr1"
	"github.com/rwtodd/Go.Sed/sed"
	"jaeger-tenant/pkg/structures"
	"strconv"
	"strings"
)

// GetValuesYaml return YAML string according to customer
func GetValuesYaml(tenant *structures.Tenant, payload *structures.TenantPayload) (string, error) {

	values, err := Asset("values-template.yaml")

	expression := "s/\\$TEMP/" + tenant.Customer + "/g; " +
		"s/\\$CASSANDRA_HOST/" + payload.CassandraHost + "/g; " +
		"s/\\$CASSANDRA_DC/" + payload.CassandraDataCenter + "/g; " +
		"s/\\$KAFKA_BROKER/" + payload.KafkaBroker + "/g; " +
		"s/\\$COLLECTOR_URL/" + tenant.JaegerCollectorURL + "/g; " +
		"s/\\$STORAGE/" + payload.StorageType + "/g; " +
		"s/\\$ELASTICSEARCH_HOST/" + payload.ESHost + "/g; " +
		"s/\\$INGESTER_ENABLED/" + strconv.FormatBool(payload.KafkaEnabled) + "/g"

	engine, _ := sed.New(strings.NewReader(expression))

	helmValue, _ := engine.RunString(string(values))

	return helmValue, err

}

// GetIngressYaml return YAML string according to customer
func GetIngressYaml(tenant *structures.Tenant, authType string) (string, error) {

	values, _ := Asset("ingress-template.yaml")

	password, _ := apr1.Hash(tenant.Password, "")

	expression := "s/\\$TEMP/" + tenant.Customer + "/g; " +
		"s/\\$QUERY_URL/" + tenant.JaegerQueryURL + "/g; " +
		"s/\\$COLLECTOR_URL/" + tenant.JaegerCollectorURL + "/g; " +
		"s/\\$AUTH_TYPE/" + authType + "/g; " +
		"s/\\$USERNAME/" + Encode(tenant.Username) + "/g; " +
		"s/\\$PASSWORD/" + Encode(tenant.Password) + "/g; " +
		"s/\\$CA_CRT/" + Encode(tenant.CACrt) + "/g; " +
		"s/\\$CA_KEY/" + Encode(tenant.CAKey) + "/g; " +
		"s/\\$TLS_CRT/" + Encode(tenant.ClientCrt) + "/g; " +
		"s/\\$TLS_KEY/" + Encode(tenant.ClientKey) + "/g; " +
		"s/\\$BASIC_AUTH/" + Encode(tenant.Username+":"+password) + "/g; "
	engine, err := sed.New(strings.NewReader(expression))
	if err != nil {
		return "", err
	}

	return engine.RunString(string(values))
}

// GetHotrodYaml return hotrod yaml according to tenant
func GetHotrodYaml(tenant *structures.Tenant) (string, error) {
	values, _ := Asset("hotrod-template.yaml")
	expression := "s/\\$TEMP/" + tenant.Customer + "/g; " +
		"s/\\$COLLECTOR_URL/" + tenant.JaegerCollectorURL + ":443/g; "
	engine, err := sed.New(strings.NewReader(expression))
	if err != nil {
		return "", err
	}
	return engine.RunString(string(values))
}

// GetPerformanceYaml returns perf yaml with env
func GetPerformanceYaml(payload *structures.PerformancePayload) (string, error) {
	values, _ := Asset("performance.yaml")

	customers := ""
	for _, name := range payload.Customers {
		customers += name + ", "
	}
	customers = customers[:len(customers)-2]

	expression := "s/\\$TEMP/" + customers + "/g; " +
		"s/\\$DOMAIN/" + payload.Domain + "/g; "
	engine, err := sed.New(strings.NewReader(expression))
	if err != nil {
		return "", err
	}
	return engine.RunString(string(values))
}
