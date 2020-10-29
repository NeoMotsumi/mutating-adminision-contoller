
#!/usr/bin/env bash

#remove old secrets/deployment/admission config 
kubectl delete deployment pod-mutation-deployment -n webhooks
kubectl delete secret pod-mutation-webhook-certs -n webhooks
kubectl delete MutatingWebhookConfiguration webhooks

docker build -t webhook-mutation:latest .
docker tag webhook-mutation:latest localhost:5000/webhook-mutation
docker push localhost:5000/webhook-mutation


tmpdir=$(mktemp -d)

basedir="$(dirname "$0")/kubernetes"

# Generate the CA cert and private key
openssl req -nodes -new -x509 -keyout ${tmpdir}/ca.key -out ${tmpdir}/ca.crt -subj "/CN=Mutating webhook CA"
# Generate the private key for the webhook server
openssl genrsa -out ${tmpdir}/webhook-server-tls.key 2048
# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -key ${tmpdir}/webhook-server-tls.key -subj "/CN=pod-mutation-service.webhooks.svc" \
    | openssl x509 -req -CA ${tmpdir}/ca.crt -CAkey ${tmpdir}/ca.key -CAcreateserial -out ${tmpdir}/webhook-server-tls.crt


# create the secret with CA cert and server cert/key
kubectl create secret generic pod-mutation-webhook-certs \
        --from-file=key.pem=${tmpdir}/webhook-server-tls.key \
        --from-file=cert.pem=${tmpdir}/webhook-server-tls.crt \
        --dry-run=client -o yaml |
    kubectl -n webhooks create -f -

kubectl apply -f ${basedir}/deployment.yaml

ca_pem_b64="$(openssl base64 -A <"${tmpdir}/ca.crt")"
sed -e 's@${CA_PEM_B64}@'"$ca_pem_b64"'@g' <"${basedir}/webhook.yaml" \
    | kubectl create -f -

rm -rf "$keydir"

echo "The webhook server has been deployed and configured!"