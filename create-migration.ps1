param(
    [string]$service,
    [string]$reason
)

migrate create -ext sql -dir services/$service/db -seq $reason
