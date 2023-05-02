# pem-energy-balance

This repository proposes a technology stack that should theoretically be able to aid in supporting the danish electrical infrastructure through a series of batteries communicating to charge / discharge power into the grid

Battery-sim
Highly concurrent simulations of batteries running in the ultra-fast golang runtime

Coordination
Kafka highly distributed fault-tolerant message broker feeding millinos of messages to a digestion pipeline (TBD in terms of tech)

Vizualisation
TBD, real-time vizualisation

## Profiles

Run docker compose with --profile=no-battery to boot without the battery-sim
Run docker compose with --profile=no-coordinator to boot without the coordinator

## Useful commands

run e.g. `docker exec -it <kafka-container-name> kafka-console-consumer --bootstrap-server localhost:9092 --topic <my-topic> --from-beginning`
to see topic messages
