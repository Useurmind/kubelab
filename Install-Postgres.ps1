helm install kubelab-postgres bitnami/postgresql
kubectl expose pod kubelab-postgres-postgresql-0 --type=NodePort --name=kubelab-postgres-node