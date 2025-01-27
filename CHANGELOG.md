## 0.2.20250127

### Game Features

* Added the Jackhammer
  * Dig down through bedrock
* Added the Flamethrower
  * Dig through up to two blocks in front of the player
  * Kills enemies and players
  * lights bombs
  * Can change number of uses (0 means unlimited)
* Added Closing Blocks
  * When a player moves through the block, the block closes after a brief pause
* Added Hideouts
  * A player can hide in a hideout and can't be hurt or chased by enemies
* Updated Demon pathfinding to use Lode Runner algorithm
* Demons can now use bars
* Changed Jumping to an item based ability (Jump Boots)
* Changed Brown player to Orange
* Added two shaders for world feel (Tab to cycle through them)
  * Underwater
  * Heat

### Editor Features

* Added a way to load Dialogs from JSON
* Menuing with arrow keys added
* Added Flamethrower options
  * Combined Item options into one
* Added palette tool
  * Used to change the color of gems, keys, and doors
  * Used to change which player can use/pickup gems, keys, tools, and doors
  * Added player colors to all items
* Added the timer value of items to the wrench display

### Bug Fixes

* Title dialog doesn't show when the title is empty
* Title dialog goes away after testing a level
* Bombs affected other bombs through bedrock
* Inventory dialogs fixed
* Cursor shows up when window isn't focused
* Players don't drop off ladder if it is 1 above bottom of screen
* Fly death animation missing when attacking player
* Player's can no longer collect or pick up other player's gems/tools
* Fixed the Floating Text Dialog
  * Color selectors weren't working
  * Leaving the input blank now removes the dialog
  * Text wasn't showing when moving between puzzles after the first time

## 0.2.20241025

### Game Features

* Added the Demon Disguise
  * When in use, enemies don't chase you or hurt you
  * Customize the time to use, the regen delay, and if it regenerates
* Updated ladder climbing animation
* In Game Title UI added
* Gamepad support added in game
* Inventory UI added
* Death count is kept track of
* Time played is kept track of
* Score is kept track of
* Async puzzle loading
* Pause menu
* Favorite levels list

### Editor Features

* ui.Dialog update, allowing for JSON based dialogs
* Added some scroll element functionality
* Added Rearrange Puzzle Set Dialog
* Created a Puzzle Preview
* Added Combine Sets Dialog
* Added a Puzzle Settings Dialog (incomplete)

### Bug Fixes

* Selection and Undo/Redo didn't reflect Wire tool changes or Timer in Metadata
* Enemies pushed each other too much, and also too little
* Multiple inputs in Dialogs could be focused at once
* Fixed Input elements and pixel.Text elements
* Exit Ladder over turf wasn't solid before exit appeared
* Crush animation was too short
* Updated dialog system to remove lag