# Project instructions

Black Sun Rising is a RTS game that's writen in Go language and the Ebiten engine.

## Behaviour

- Don't invent additional game mechanics, lore or story parts without aggrements.
- Read design documents before implementing a specific feature.

## Code style

- Define interfaces in the place of using.
- Avoid using log package, use slog instead.
- Don't use os.Exit outside main function.

## Architecture

- Game design documents are lying in the design folder.
- Separate rendering and game logic.
- Prefer straight-forward logic to complex OOP or ECS solutions.
- Follow old 80s style of the software architecture.
