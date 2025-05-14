# Create the raycast dummy to track maximum attemps when raycasting
scoreboard objectives add copper_items.ray_steps dummy "Raycast Max Steps"

# Create a global on/off for debug
scoreboard objectives add copper_items.debug dummy "Enable Debug Messages"
scoreboard players set #global copper_items.debug 0

tellraw @a "Datapack reloaded"