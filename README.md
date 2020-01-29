# Jaeger Tenant Deploy

Kubernetes Job to deploy Jaeger components

To get started, grab the `jaeger-tenant.yaml` and change the environment variable according to your need, then

```bash
$ kubectl apply -f jaeger-tenant.yaml
```

### Environment variables


| Variable              | Use case |
| --------------------- | -------- |
| CUSTOMER_NAME         | Tenant name must be unique  |
| KAFKA_BROKER          | Comma separated list of Kafka brokers with port number  |
| CASSANDRA_HOST        | Hostname of Cassandra DB  |
| DOMAIN                | Domain of ingress, used for creating sub-domain for collector and query |
| CASSANDRA_DATACENTER  | Datacenter for cassandra  |
| KUBECONFIG            | Path to Kube config file inside container |

### Volume Mount

A configmap is required consisting of Kube config which will be required by container to deploy components on a K8s cluster.

Create the configmap prior to deploying this Job

```bash
$ kubectl create configmap kube-config --from-file=kube.yaml=<path to config in local system>
```

This configmap will get mounted in container.
