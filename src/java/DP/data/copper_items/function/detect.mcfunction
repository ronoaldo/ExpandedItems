tellraw @a "Procurando Copper Stick"

execute as @e[type=minecraft:item] \
    if items entity @s contents minecraft:stick[minecraft:custom_model_data={floats:[42006]}] at @s \
    if block ~ ~-0.01 ~ minecraft:crafting_table \
    run return run tellraw @a "Achou"

tellraw @a "Não Achou"