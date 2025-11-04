# Black Sun Rising - Main Design Document

## Overview

_Black Sun Rising_ is a single-player real-time strategy game focused on storytelling rather than competitive ladder play. Set in a hard sci-fi universe without magic, psychic powers, or other fantastical elements, the game follows realistic physics and technology principles. The game is built with Go and the Ebiten engine. All battles take place in space, featuring stars, planets, and asteroids as the primary environment.

Unlike traditional RTS games, there are no static buildings; instead, all economic activities are managed by large ships such as Motherships and Cruisers. Interstellar communication and transportation are achieved through realistic methods such as wormholes and Alcubierre drives.

## Core Systems

### Game World

The game world consists entirely of space environments including:
- Stars.
- Planets of various sizes and types.
- Asteroid fields.
- Debris from destroyed or abandoned ships.

All gameplay takes place in these space environments, with no planetary surfaces or ground-based operations.

### Units and Buildings

Instead of traditional buildings, the game features specialized ships that serve functional roles:

**Capital Ships:**
- Mothership: The primary command vessel and center of operations
- Cruisers: Large support ships that assist with various functions

**Specialized Ships:**
- Worker Ships: Collect resources from asteroids and ship wreckage
- Combat Ships: Various classes for offensive and defensive operations
- Support Ships: Provide auxiliary functions like repairs and logistics

All economic and production activities occur aboard these mobile platforms rather than fixed structures.

### Resources and Economy

The game features a simplified economy based on a single resource type:

**Materials:** The primary and only resource in the game, used for all production and upgrades.

**Collection Methods:**
- Specialized worker ships harvest materials from asteroid fields
- Salvage operations recover materials from crashed or abandoned ships
- Resource collection is an active process requiring ship deployment and management

All resource processing and storage occurs aboard the fleet's capital ships, making fleet positioning and protection critical for economic success.

### Combat System

Combat takes place entirely in space between fleets of ships:

- Ship-to-ship engagements using various weapon systems
- Fleet composition and positioning are crucial for victory
- Damage to ships affects their performance and may destroy them
- No ground combat or planetary invasions

### User Interface

### Story Mode

As a single-player focused game, the emphasis is on narrative-driven missions and campaign progression rather than competitive multiplayer:

- Campaign mode with story progression
- Mission-based objectives tied to narrative events
- Character development and dialogue systems
- Branching storylines based on player choices and performance
