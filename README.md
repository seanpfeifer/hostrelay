# Host Relay

The point of this is to be a dedicated relay server for multiplayer applications with two needs:

* Reliable, ordered data transmission (TCP-like)
* Unreliable data transmission (UDP)

This exists to allow players to join games that don't need dedicated server logic.

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
