#!/bin/bash

# Initializing variables
customerName=$CUSTOMER_NAME
domain=$DOMAIN
collectorURL=$CUSTOMER_NAME"-collector."$DOMAIN
queryURL=$CUSTOMER_NAME"-query."$DOMAIN
kafkaBroker=$KAFKA_BROKER
cassandraHost=$CASSANDRA_HOST
cassandraDC=$CASSANDRA_DATACENTER

mkdir $customerName 

# Creating openssl config
cat > $customerName/csr_details.txt <<-EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[ dn ]
C=IN
ST=Bangalore
L=Bangalore
O=SAP Labs
OU=Jaeger Service Provider
emailAddress=arpit.agarwal02@sap.com
CN = $customerName

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $collectorURL
EOF

# Generating CA.key and CA.crt
openssl req -new -x509 -sha256 -newkey rsa:2048 -nodes -keyout $customerName/CA.key -days 365 -out $customerName/CA.crt -config $customerName/csr_details.txt

# Generating KEY and CSR
openssl req -nodes -newkey rsa:2048 -keyout $customerName/server.key -out $customerName/server.csr -config $customerName/csr_details.txt

# Generating CRT
# Signing CSR with CA.crt and CA.key
openssl x509 -req -days 1460 -in $customerName/server.csr -CA $customerName/CA.crt  -CAkey $customerName/CA.key -set_serial 01 -out $customerName/server.crt

# Updating repo to latest version
helm repo update

:'
# Generating Jaeger components YAML
helm template --set \
provisionDataStore.cassandra=false,\
storage.cassandra.keyspace=$customerName,\
storage.cassandra.host=$cassandraHost,\
cassandra.config.dc_name=$cassandraDC,\
storage.cassandra.password=password,\
ingester.enabled=true,\
storage.kafka.brokers={$kafkaBroker},\
storage.kafka.topic=$customerName,\
collector.cmdlineParams.collector_grpc_tls=false,\
collector.cmdlineParams.collector_grpc_tls_cert=/tls/server.crt,\
collector.cmdlineParams.collector_grpc_tls_client-ca=/tls/ca.crt,\
collector.cmdlineParams.collector_grpc_tls_key=/tls/server.key,\
collector."extraSecretMounts[0]".name=jaeger-tls,\
collector."extraSecretMounts[0]".mountPath=/tls,\
collector."extraSecretMounts[0]".readOnly=true,\
collector."extraSecretMounts[0]".secretName=$customerName-tls-config,\
agent.enabled=false,\
agent.cmdlineParams."reporter\.grpc\.host-port"=$collectorURL:443,\
agent.cmdlineParams.reporter_grpc_tls=true,\
agent.cmdlineParams.reporter_grpc_tls_ca=/tls/ca.crt,\
agent.cmdlineParams.reporter_grpc_tls_cert=/tls/tls.crt,\
agent.cmdlineParams.reporter_grpc_tls_key=/tls/tls.key,\
agent."extraSecretMounts[0]".name=jaeger-tls,\
agent."extraSecretMounts[0]".mountPath=/tls,\
agent."extraSecretMounts[0]".readOnly=true,\
agent."extraSecretMounts[0]".secretName=$customerName-tls-config \
$customerName jaegertracing/jaeger > $customerName/jaeger.yaml
'

# Generating Ingress YAML
sed \
-e 's/temp/'$customerName'/g' \
-e 's/query-url/'$queryURL'/g' \
-e 's/collector-url/'$collectorURL'/g' \
ingress-template.yaml > $customerName/ingress.yaml

# Generating Values YAML
sed \
-e 's/temp/'$customerName'/g' \
-e 's/cassandra-host/'$cassandraHost'/g' \
-e 's/cassandra-dc/'$cassandraDC'/g' \
-e 's/kafka-broker/'$kafkaBroker'/g' \
-e 's/collector-url/'$collectorURL'/g' \
values-template.yaml > $customerName/values.yaml

# Generating Jaeger YAML
helm template $customerName . \
--values $customerName/values.yaml > $customerName/jaeger.yaml

# Creating Secrets and applying YAMLs
kubectl create secret generic "$customerName-tls-config" \
--from-file=tls.crt=$customerName/server.crt \
--from-file=tls.key=$customerName/server.key \
--from-file=ca.crt=$customerName/CA.crt \
--from-file=ca.key=$customerName/CA.key 

kubectl apply -f $customerName/jaeger.yaml $customerName/ingress.yaml
