
customerName=$CUSTOMER_NAME

tls_ca=`kubectl get secrets $customerName-tls-config -o jsonpath='{.data.ca\.crt}' | base64 --decode`
tls_ca_key=`kubectl get secrets $customerName-tls-config -o jsonpath='{.data.ca\.key}' | base64 --decode`
tls_cert=`kubectl get secrets $customerName-tls-config -o jsonpath='{.data.tls\.crt}' | base64 --decode`
tls_key=`kubectl get secrets $customerName-tls-config -o jsonpath='{.data.tls\.key}' | base64 --decode`

collectorURL=`kubectl get ingress $customerName-collector | awk '{print $2}' | grep -v HOSTS`
collectorURL="$collectorURL:443"
queryURL=`kubectl get ingress $customerName-query | awk '{print $2}' | grep -v HOSTS`
queryURL="https://$queryURL"

echo
echo "JSON to be used in VCAP"
echo
jq -n   --arg tls_ca "$tls_ca" \
        --arg tls_ca_key "$tls_ca_key" \
        --arg tls_cert "$tls_cert" \
        --arg tls_key "$tls_key" \
        --arg collector "$collectorURL" \
        --arg query "$queryURL" \
        '{"jaeger-collector-url": $collector, "jaeger-ui-url": $query, "tls_ca": $tls_ca, "tls_cert": $tls_cert, "tls_key": $tls_key, "tls_ca_key": $tls_ca_key}'

echo 