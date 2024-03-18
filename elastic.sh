#!/bin/bash

INDEX=cars

curl -X GET "localhost:9200/${INDEX}/_mapping?pretty"
curl -X GET "localhost:9200/${INDEX}/_settings?pretty"
