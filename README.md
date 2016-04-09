# Empire Destiny / Gods and Empires #

# Flow #

Match is created

One player is selected to start

PlayerA starts

## Round 1
### Player A Turn
PlayerA -> Harvest()
PlayerA -> DoActions() # attach followers and items
|
if (palyerA attack)
  PlayerA -> assign unit
  PlayerB -> assign unit
  loop:
    PlayerA -> Action()
    PlayerB -> Action()
  --> Resolve battle # settle life force between sides
PlayerA -> Build() # buy holdings and personalities
### Player B Turn
PlayerB -> Harvest()
PlayerB -> DoActions()
|
if (playerB attack)
  PlayerB -> assign unit
  PlayerA -> assign unit
  loop:
    PlayerB -> Action()
    PlayerA -> Action()
  --> Resolve battle # settle life force between sides
# Traits #

## Rules ##

### Context ###
playerA - Attacking Player
playerB - Defending Player
card in action
cards affected by action

each exp
  Find all cards that match 'target' of trait exp
  for each card
    apply action
