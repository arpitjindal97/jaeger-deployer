# Jaeger Deployer

Application to deploy Jaeger components

## Prerequisite

Cassandra, ElasticSearch and Kafka should be installed prior before you go for installation of Tenant via this application

Below are the commands through which you can install above backing stateful sets :

### Installation of Cassandra

```bash
helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
helm install cassandra incubator/cassandra --values https://github.wdf.sap.corp/raw/EngSrv-Monitoring/jaeger-deployer/master/cassandra-values.yaml
```

### Installation of Kafka

```bash
helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
helm install kafka incubator/kafka --values https://github.wdf.sap.corp/raw/EngSrv-Monitoring/jaeger-deployer/master/kafka-values.yaml
```

### ConfigMap with KubeConfig

```bash
kubectl create configmap kube-config --from-file=config=<path-to-file>
```

## Deploy

```bash
kubectl apply -f https://github.wdf.sap.corp/raw/EngSrv-Monitoring/jaeger-deployer/master/jaeger-tenant.yaml
```

## Usage

#### Creating Tenants
```
PUT /tenant HTTP/1.1
Content-Type: application/json

{
	"customers": [
		"cust1",
		"cust2"
		],
	"domain": "ingress.mon-team.perf-load.shoot.dev.k8s-hana.ondemand.com",
	"storageType": "elasticsearch",
	"esHost":"elasticsearch-master",
	"cassandraHost":"cassandra",
	"cassandraDatacenter":"dc1",
	"kafkaEnabled":false,
	"kafkaBroker": "jaeger-kafka:9092",
	"hotrodExample": true,
	"authType": ""
}

```

#### Getting Tenant Detail
```
GET /tenant/{customerName} HTTP/1.1
```

#### Deleting Tenant
```
DELTE /tenant/{customerName} HTTP/1.1
```

#### Deploying Performance Job
```
PUT /performance
Content-Type: application/json

{
	"customers": [
		"cust1",
		"cust2"
		],
	"domain": "ingress.mon-team.perf-load.shoot.dev.k8s-hana.ondemand.com",
    "threadNumber": "100"
}
```
