
customerName="cust"

kubectl delete \
serviceaccount/$customerName-jaeger-cassandra-schema \
serviceaccount/$customerName-jaeger-collector \
serviceaccount/$customerName-jaeger-query \
serviceaccount/$customerName-jaeger-agent \
serviceaccount/$customerName-jaeger-ingester \
service/$customerName-jaeger-agent \
service/$customerName-jaeger-collector \
service/$customerName-jaeger-query \
daemonset.apps/$customerName-jaeger-agent \
deployment.apps/$customerName-jaeger-collector \
deployment.apps/$customerName-jaeger-ingester \
deployment.apps/$customerName-jaeger-query \
job.batch/$customerName-jaeger-cassandra-schema \
secrets/$customerName-jaeger-cassandra

kubectl delete secret $customerName-tls-config $customerName-tls-query

kubectl delete ingress $customerName-collector $customerName-query
