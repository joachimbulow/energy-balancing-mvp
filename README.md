# pem-energy-balance

This repository proposes a technology stack that should theoretically be able to aid in supporting the danish electrical infrastructure through a series of batteries communicating to charge / discharge power into the grid

It simulates the protocol: Packetized Energy Management, which proposes breaking energy into packets for coordination

This setup simulaltes this deployment at state-grid scale with 1 million+ batteries running simultaneously

![image](https://github.com/joachimbulow/energy-balancing-mvp/assets/43139346/3bd6cf94-3eb2-4781-bb07-bb953ce1b9f4)


## Components 

Battery-sim
Highly concurrent simulations of batteries running concurrently in the ultra-fast golang runtime

Coordination
Kafka highly distributed fault-tolerant message broker feeding millinos of messages to a digestion pipeline
written and deployed in an Apache Flink cluster

Vizualisation
InfluxDB timeseries

Transmission System Operator
Written in Nodejs

Deployed on Kubernetes with ArgoCD and Helm
https://github.com/joachimbulow/energy-balancing-mvp-deployment

## Local deployment
Run docker compose with --profile=all to boot everything
Run docker compose with --profile=no-battery to boot without the battery-sim
Run docker compose with --profile=no-coordinator to boot without the coordinator
