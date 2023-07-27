# Design note

## Overview

Monitors:

- iface monitor
- socket monitor

Retreive metrics from those monitors, and put them into datastore.

## Data format

We should support compressing, storage quotas and do rotating.

Collect raw metrics, do some simple stats to get some events.

For a simple choice, use BSON and BadgerDB.

### Record format

- Timestamp. eg. 2023-07-27 12:00:00.123
- Type. iface, socket, ...
- Epoch. Auto inc through running
- Payload. The metric data here.

## Export metrics

Export the database into a single file, expose them through HTTP. Then we can use a simple script to download them and
get them back.
