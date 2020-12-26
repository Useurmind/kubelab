param(
    $dbName = "projects"
)


$password = (kubectl get secrets -n kubelab kubelab-postgres-postgresql -o jsonpath="{.data.postgresql-password}" | gobase64.exe -d).Trim()

$env:KUBELAB_DB_HOST = "localhost"
$env:KUBELAB_DB_PORT = "30432"
$env:KUBELAB_DB_DBNAME = $dbName
$env:KUBELAB_DB_USER = "postgres"
$env:KUBELAB_DB_PASSWORD = $password