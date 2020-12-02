# Host Relay Prototype

The point of this is to be a dedicated relay server for multiplayer applications with two needs:

* Reliable, ordered data transmission (TCP-like)
* Unreliable data transmission (UDP)

This exists to allow players to join games that don't need dedicated server logic.

**DISCLAIMER: This is a hacked together prototype and should not be used in production! Code will change without warning and nothing is guaranteed!**

## What this **doesn't have**

* Anti-cheat
* Data sanity checking
* Anti-spam/congestion control
* Complex authentication

## Basic design

* TCP socket listening for connections
  * Each connection is stored
* UDP socket to listen for messages from each player
  * Broadcast each received message to the rest of the connected players

## Building & running

```sh
go run .\cmd\server
```

To run a test client in TCP mode:

```sh
go run .\cmd\client
# Input strings and press enter to send. Type `exit` to quit.
```

For UDP:

```sh
go run .\cmd\client -udp
# Input strings and press enter to send. Type `exit` to quit.
```
