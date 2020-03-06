package kubectl

import (
	"bytes"
	"errors"
	"jaeger-tenant/pkg/structures"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os/exec"
	"strings"
)

// GetTenant returns Tenant struct fetching details from K8s
func GetTenant(customer string) (tenant *structures.Tenant, err error) {

	context := NewContext()
	context.Make()

	ingressCollector, err := context.GetClient().ExtensionsV1beta1().Ingresses(namespace).Get(customer+"-collector", metav1.GetOptions{})
	ingressQuery, err := context.GetClient().ExtensionsV1beta1().Ingresses(namespace).Get(customer+"-query", metav1.GetOptions{})
	secret, err := context.GetClient().CoreV1().Secrets(namespace).Get(customer+"-tls", metav1.GetOptions{})

	if err == nil {

		tenant = &structures.Tenant{
			Customer:           customer,
			JaegerCollectorURL: ingressCollector.Spec.Rules[0].Host + ":443",
			JaegerQueryURL:     "https://" + ingressQuery.Spec.Rules[0].Host,
			Username:           string(secret.Data["username"]),
			Password:           string(secret.Data["password"]),
			CACrt:              string(secret.Data["ca.crt"]),
			CAKey:              string(secret.Data["ca.key"]),
			ClientCrt:          string(secret.Data["tls.crt"]),
			ClientKey:          string(secret.Data["tls.key"]),
		}
	}
	return
}

// ApplyYaml applies YAML to K8s using `kubectl apply`
func ApplyYaml(tenantYaml string) (err error) {

	cmd := exec.Command("bash", "-c", "cat << EOF | kubectl apply -f -\n"+tenantYaml+"\nEOF")

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = errors.New(err.Error() + "\n" + stderr.String())
	}
	return
}

// DeleteTenant delete all K8s resources of given tenant
func DeleteTenant(customer string) (respError error) {

	context := NewContext()
	context.Make()

	respError = errors.New("no resources found related to " + customer)

	err1 := context.GetClient().AppsV1().Deployments(namespace).Delete(customer+"-jaeger-collector", &metav1.DeleteOptions{})
	err2 := context.GetClient().AppsV1().Deployments(namespace).Delete(customer+"-jaeger-query", &metav1.DeleteOptions{})
	err3 := context.GetClient().AppsV1().Deployments(namespace).Delete(customer+"-jaeger-ingester", &metav1.DeleteOptions{})

	err4 := context.GetClient().BatchV1().Jobs(namespace).Delete(customer+"-jaeger-cassandra-schema", &metav1.DeleteOptions{})

	err5 := context.GetClient().CoreV1().Secrets(namespace).Delete(customer+"-jaeger-cassandra", &metav1.DeleteOptions{})
	err6 := context.GetClient().CoreV1().Secrets(namespace).Delete(customer+"-tls", &metav1.DeleteOptions{})
	err7 := context.GetClient().CoreV1().Secrets(namespace).Delete(customer+"-tls-query", &metav1.DeleteOptions{})

	err8 := context.GetClient().CoreV1().Services(namespace).Delete(customer+"-jaeger-agent", &metav1.DeleteOptions{})
	err9 := context.GetClient().CoreV1().Services(namespace).Delete(customer+"-jaeger-collector", &metav1.DeleteOptions{})
	err10 := context.GetClient().CoreV1().Services(namespace).Delete(customer+"-jaeger-query", &metav1.DeleteOptions{})

	err11 := context.GetClient().CoreV1().ServiceAccounts(namespace).Delete(customer+"-jaeger-cassandra-schema", &metav1.DeleteOptions{})
	err12 := context.GetClient().CoreV1().ServiceAccounts(namespace).Delete(customer+"-jaeger-collector", &metav1.DeleteOptions{})
	err13 := context.GetClient().CoreV1().ServiceAccounts(namespace).Delete(customer+"-jaeger-query", &metav1.DeleteOptions{})
	err14 := context.GetClient().CoreV1().ServiceAccounts(namespace).Delete(customer+"-jaeger-agent", &metav1.DeleteOptions{})
	err15 := context.GetClient().CoreV1().ServiceAccounts(namespace).Delete(customer+"-jaeger-ingester", &metav1.DeleteOptions{})

	err16 := context.GetClient().ExtensionsV1beta1().Ingresses(namespace).Delete(customer+"-collector", &metav1.DeleteOptions{})
	err17 := context.GetClient().ExtensionsV1beta1().Ingresses(namespace).Delete(customer+"-query", &metav1.DeleteOptions{})

	err18 := context.GetClient().AppsV1().Deployments(namespace).Delete(customer+"-hotrod", &metav1.DeleteOptions{})
	err19 := context.GetClient().CoreV1().Services(namespace).Delete(customer+"-hotrod", &metav1.DeleteOptions{})

	list, _ := context.GetClient().CoreV1().Pods(namespace).List(metav1.ListOptions{})
	err20 := errors.New("no pod found")
	for _, pod := range list.Items {
		if strings.Contains(pod.Name, customer+"-jaeger") {
			err20 = context.GetClient().CoreV1().Pods(namespace).Delete(pod.Name, &metav1.DeleteOptions{})
		}
	}

	// if anyone one of them is nil, then something is deleted
	if err1 == nil || err2 == nil || err3 == nil || err4 == nil || err5 == nil || err6 == nil || err7 == nil || err8 == nil ||
		err9 == nil || err10 == nil || err11 == nil || err12 == nil || err13 == nil || err14 == nil || err15 == nil || err16 == nil ||
		err17 == nil || err18 == nil || err19 == nil || err20 == nil {
		respError = nil
	}

	return
}
