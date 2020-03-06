package jaeger

import (
	"encoding/json"
	"jaeger-tenant/pkg/certificates"
	"jaeger-tenant/pkg/helm"
	"jaeger-tenant/pkg/kubectl"
	"jaeger-tenant/pkg/structures"
	"jaeger-tenant/utils"
)

// CreateTenant tries to create tenant and returns message
func CreateTenant(customer string, payload *structures.TenantPayload) (tenant structures.Tenant) {

	tenant = structures.Tenant{
		Customer:           customer,
		JaegerCollectorURL: "",
		JaegerQueryURL:     "",
		Username:           "",
		Password:           "",
		CACrt:              "",
		CAKey:              "",
		ClientCrt:          "",
		ClientKey:          "",
		Error:              "",
	}
	_, err := kubectl.GetTenant(customer)
	if err == nil {
		tenant.Error = "Tenant already exists"
		return
	}
	tenant.JaegerCollectorURL = customer + "-collector." + payload.Domain
	tenant.JaegerQueryURL = customer + "." + payload.Domain

	certificates.GenerateCredentials(&tenant)
	tenantYaml, err := helm.GetTenantYaml(&tenant, payload)
	if err != nil {
		tenant = structures.Tenant{
			Error: err.Error(),
		}
		return
	}
	ingress, _ := utils.GetIngressYaml(&tenant, payload.AuthType)
	if payload.HotrodExample {
		hotrod, _ := utils.GetHotrodYaml(&tenant)
		ingress = ingress + hotrod
	}
	tenantYaml = ingress + tenantYaml

	err = kubectl.ApplyYaml(tenantYaml)
	if err != nil {
		tenant = structures.Tenant{
			Error: err.Error(),
		}
		return
	}

	tenant.JaegerCollectorURL += ":443"
	tenant.JaegerQueryURL = "https://" + tenant.JaegerQueryURL
	return
}

// GetTenant return tenant information from K8s
func GetTenant(customer string) []byte {

	context := kubectl.NewContext()
	context.Make()

	tenant, err := kubectl.GetTenant(customer)
	var b []byte
	if err != nil {
		resp := utils.CreateErrorResponse(err)
		b, _ = json.Marshal(resp)
	} else {
		b, _ = json.Marshal(tenant)
	}

	return b
}

// DeleteTenant deletes the tenant, returning a message
func DeleteTenant(customer string) []byte {

	b, _ := json.Marshal(utils.CreateMessageResponse("successfully removed all possible resources"))

	err := kubectl.DeleteTenant(customer)
	if err != nil {
		b, _ = json.Marshal(utils.CreateErrorResponse(err))
	}

	return b
}

// CreatePerformanceTest deploys performance test for a bunch of tenant
func CreatePerformanceTest(payload *structures.PerformancePayload) error {

	perfYaml, err := utils.GetPerformanceYaml(payload)
	if err != nil {
		return err
	}
	return kubectl.ApplyYaml(perfYaml)
}
