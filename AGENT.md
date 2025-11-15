# Project instructions

Black Sun Rising is a RTS game that's writen in Go language and the Ebiten engine.

## Behaviour

- Don't invent additional game mechanics, lore or story parts without aggrements.
- Read design documents before implementing a specific feature.

## Code style

- Define interfaces in the place of using.
- Avoid using log package, use slog instead.
- Don't use os.Exit outside the main function.
- Avoid magic numbers and declare constants inside const or var blocks.

## Architecture

- Game design documents are lying in the design folder.
- Separate game logic, input and rendering.
- Preserve all state that is needed for all components in the scene.State.
- Store individual logic of game objects inside it's methods and interaction logic in separated services.
