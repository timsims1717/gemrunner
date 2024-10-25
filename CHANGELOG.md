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

### Bugs

* Selection and Undo/Redo didn't reflect Wire tool changes or Timer in Metadata
* Enemies pushed each other too much, and also too little
* Multiple inputs in Dialogs could be focused at once
* Fixed Input elements and pixel.Text elements
* Exit Ladder over turf wasn't solid before exit appeared
* Crush animation was too short
* Updated dialog system to remove lag