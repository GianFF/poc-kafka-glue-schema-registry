#!/bin/bash
set -e
cd java
mvn clean install
mvn clean package
java -jar target/kafka-glue-poc-1.0.0.jar
