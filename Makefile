NAME=ExpandedItems

DP=output/datapacks/$(NAME)-datapack.zip
RP=output/resourcepacks/$(NAME)-resourcepack.zip
MP=output/mcpacks/$(NAME)-resourcepack.mcpack
GM=output/mcpacks/$(NAME)-geyser_mappings.json

all: package

package:
	go run cmd/creeper/main.go -n $(NAME) -i expanded_items.png
	cp src/bedrock/geyser_mappings.json output/mcpacks/$(NAME)-geyser_mappings.json

server: package
	mkdir -p output/server/world/datapacks/ \
		output/server/plugins/Geyser-Spigot/packs/ \
		output/server/plugins/Geyser-Spigot/custom_mappings/
	cp $(DP) output/server/world/datapacks/
	cp $(MP) output/server/plugins/Geyser-Spigot/packs/
	cp $(GM) output/server/plugins/Geyser-Spigot/custom_mappings/
	if [ -f docker-compose.override.yaml ] ; then \
		sed -e "s/SHA1: .*/SHA1: $$(sha1sum < $(RP) | cut -f 1 -d ' ')/" \
			-i docker-compose.override.yaml ; \
	fi
	docker compose up

clean:
	rm -rf output

vanilla:
	mkdir -p vanilla
	curl https://raw.githubusercontent.com/misode/mcmeta/refs/heads/summary/item_components/data.json > vanilla/item_components.json