# Simulation

This directory contains the simulation code for the project. The simulation is written in
Python and will send the data to `ws://localhost:4510` by default. The simulation sends
data with some random information to simulate the real world in a given time interval.

The data being sent is in the following format:

```json
{
  "year": 2020,
  "month": 1,
  "day": 30,
  "hour": 8,
  "utc_timestamp": "2020-01-30T16:00:00.000Z",
  "start_junction_id": "8aaf6ad421333ad741cb1d5de1e3fa83e7d0e908",
  "end_junction_id": "8bbd97259361a23d374e60e6582019550fc58e0f",
  "osm_way_id": 417094233,
  "osm_start_node_id": 4714793573,
  "osm_end_node_id": 1014244233,
  "speed_mph_mean": 27.761,
  "speed_mph_stddev": 3.679
}
```

## Running the simulation

The simulation was developed and tested with the Python version `3.11.3`. To run the
simulation, you will need to install the dependencies with `pip install -r requirements.txt`.
After that, you can run the simulation with `python simulation.py`. The simulation will
run indefinitely until you stop it by killing the process.
