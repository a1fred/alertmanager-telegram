#! /usr/bin/env sh
# Tnx to https://blog.mafr.de/2020/09/13/testing-alertmanager/

URL="http://localhost:9093/api/v1/alerts"

curl -si -X POST -H "Content-Type: application/json" "$URL" -d '
[
  {
    "labels": {
      "alertname": "InstanceDown",
      "instance": "localhost:8080",
      "job": "node",
      "severity": "critical"
    },
    "annotations": {
      "summary": "Instance is down"
    },
    "generatorURL": "http://localhost:9090/graph"
  }
]
'
