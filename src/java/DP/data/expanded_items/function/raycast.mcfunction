# Check if the current position we are casting to has a glow_item_frame around 0.2 blocks
# If found, we tag it as the raycast_found item and stop the search
execute if entity @e[type=glow_item_frame,distance=..0.2,sort=nearest,limit=1] \
  run return run tag @e[type=glow_item_frame,distance=..0.2,sort=nearest,limit=1] add raycast_found

# Otherwise, decrease the step count ...
scoreboard players remove @s expanded_items.ray_steps 1

# ... and recurse if not reached zero
execute if score @s expanded_items.ray_steps matches 1.. \
  positioned ^ ^ ^0.1 run function expanded_items:raycast