#!/bin/bash

echo "Deploying Grafana..."
helm install grafana helm/charts/grafana --namespace monitoring --create-namespace
echo "Deploying Loki..."
helm install loki helm/charts/loki-stack --namespace monitoring
echo "Deploying MongoDB..."
helm install mongodb helm/charts/mongodb --namespace database --create-namespace
echo "Deploying RabbitMQ..."
helm install rabbitmq helm/charts/rabbitmq --namespace messaging --create-namespace
echo "Deploying api-gateway microservice"
helm install api-gateway ./microservices/api-gateway --namespace microservices --create-namespace
helm install api-gateway ./microservices/account-service --namespace microservices

