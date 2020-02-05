#!/bin/bash

# Initializing variables
customerName=$CUSTOMER_NAME
domain=$DOMAIN
collectorURL=$CUSTOMER_NAME"-collector."$DOMAIN
queryURL=$CUSTOMER_NAME"."$DOMAIN
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

echo
echo "Generating Certificates for Tenant"

# Generating CA.key and CA.crt
openssl req -new -x509 -sha256 -newkey rsa:2048 -nodes \
            -keyout $customerName/CA.key \
            -days 365 \
            -out $customerName/CA.crt \
            -config $customerName/csr_details.txt &> /dev/null

# Generating KEY and CSR
openssl req -nodes -newkey rsa:2048 \
            -keyout $customerName/server.key \
            -out $customerName/server.csr \
            -config $customerName/csr_details.txt &> /dev/null

# Generating CRT
# Signing CSR with CA.crt and CA.key
openssl x509    -req -extfile <(printf "subjectAltName=DNS:$collectorURL") \
                -days 1460 \
                -in $customerName/server.csr \
                -CA $customerName/CA.crt \
                -CAkey $customerName/CA.key \
                -set_serial 01 \
                -out $customerName/server.crt &> /dev/null

echo 
echo "Updating Helm Repository"
echo
# Updating repo to latest version
helm repo update

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

echo 
echo "Installing Jaeger Components"
echo
# Generating Jaeger YAML
helm template $customerName jaegertracing/jaeger \
--values $customerName/values.yaml > $customerName/jaeger.yaml

# Creating Secrets and applying YAMLs
kubectl create secret generic "$customerName-tls-config" \
--from-file=tls.crt=$customerName/server.crt \
--from-file=tls.key=$customerName/server.key \
--from-file=ca.crt=$customerName/CA.crt \
--from-file=ca.key=$customerName/CA.key 

kubectl apply -f $customerName/jaeger.yaml -f $customerName/ingress.yaml

./get_json.sh
