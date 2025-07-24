#!/bin/bash
set -e
docker-compose up -d
docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --create --topic users.signedup --partitions 1 --replication-factor 1 || true
docker-compose exec kafka kafka-topics.sh --bootstrap-server localhost:9092 --list
