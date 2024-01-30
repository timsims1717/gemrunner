# Controls

* Left: move left, change held item to side
* Right: move right, change held item to side
* Up: move up ladders, change held item to up
* Down: move down ladders, drop off boxes and ropes, drop (not use) carried items
* Jump: jump from the ground (or boxes)
* Pick up/Use: Pick up an item, use that item

# Jumping

* Can't jump from ladders
* Can't jump if there is a block above you
* Two kinds of jumps:
  * "long jump" that crosses 1 wide gaps
  * "high" jump that goes up one block and over in the direction facing
* If there are blocks up/right and up/left, high jump
* If there is a block left or right (when facing that way), high jump
* Otherwise, long jump

# Items that can be Carried

* Box
  * Can be walked on from above
  * Can be jumped on
  * Can be dropped down from (press down)
  * When used, just sets it down
* Key
  * When used at a locked door, unlock that door
* Regular Bombs
  * When used, three second timer, then explode in same area as LR
* Plus Bombs
  * When used, three second timer, then explode in a plus shape, 5 tiles wide and tall
* Bear Trap
  * When used, place down and set. Whenever another player or walking enemy steps there, kill them and disappear
* Snare
  * When used, place down and set. Whenever another player of walking enemy steps there, trap them, if trapped character ever leaves, disappear
* Light
  * Gives off extra light in dark levels 
  * When used, just sets it down
* Other players
  * That player can get away by jumping, can grab stuff

# Puzzle Load Dialog Information

* Name of puzzle group (filename) (only one for now)
* \# of players
* \# of puzzles
* Creator's Name
* Two side elements
  * First level layout
  * Description