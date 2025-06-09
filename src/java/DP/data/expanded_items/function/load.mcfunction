# Create the raycast dummy to track maximum attemps when raycasting
scoreboard objectives add expanded_items.ray_steps dummy "Raycast Max Steps"
scoreboard objectives add expanded_items.current_repair_cost dummy "Current Item Repair Cost"
 
# Create a global on/off for debug
scoreboard objectives add expanded_items.debug dummy "Enable Debug Messages"
scoreboard players set #global expanded_items.debug 0

tellraw @a "Datapack reloaded"