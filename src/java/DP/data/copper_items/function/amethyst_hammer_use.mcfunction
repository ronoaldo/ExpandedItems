# Revoke the advancement
advancement revoke @s only copper_items:amethyst_hammer_use

# Raycast to find the possible item frame user is looking at
execute as @s anchored eyes positioned ^ ^ ^ anchored feet \
  run function copper_items:raycast_start

# Return if no item frame was found
execute unless entity @e[tag=raycast_found] \
  run return run tellraw @s \
  "Ops... Try again but aim at the item!"

# Detect if the itemframe is on top of amethyst block.
# Otherwise, tell the player and return
execute at @e[tag=raycast_found] unless block ~ ~-1 ~ minecraft:amethyst_block \
  run tellraw @s \
  "Ops... You need an Amethyst Block for the magick to happen!"
execute at @e[tag=raycast_found] unless block ~ ~-1 ~ minecraft:amethyst_block \
  run return run function copper_items:raycast_end

# Get the current repair cost value and store it
execute store result score @e[tag=raycast_found] copper_items.current_repair_cost \
  run data get entity @e[tag=raycast_found,limit=1] Item.components.minecraft:repair_cost

# If we are all set, let's revert the repair cost back to 38
data modify entity @e[tag=raycast_found,limit=1,scores={copper_items.current_repair_cost=39..}] \
  Item.components.minecraft:repair_cost set value 38

# Show the particles to display success
execute at @e[tag=raycast_found,scores={copper_items.current_repair_cost=39..}] \
  run particle minecraft:dust{color:[1, 0, 1], scale:1} ~ ~0.3 ~ 0.125 0.125 0.125 1 10

# Cleanup
function copper_items:raycast_end
