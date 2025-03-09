NAME=ExpandedItems

DP=output/datapacks/$(NAME)-datapack.zip
RP=output/resourcepacks/$(NAME)-resourcepack.zip
MP=output/mcpacks/$(NAME)-resourcepack.mcpack
GM=output/mcpacks/$(NAME)-geyser_mappings.json

BN_0=$(shell date +%j | bc )
BN_1=$(shell echo 60*$$(date +%H)+$$(date +%M) | bc)
BN_2=$(shell echo $$(date +%S)/10 | bc)

all: package

validate:
	find src -type f -iname *.json | while read f ; do \
		echo "Checking $$f" ;\
		jq . $$f || exit 1 ;\
	done

build-number:
	sed -e "s/\"BUILD_NUMBER\"/$(BN_1)$(BN_2)/g" -e "s/\"BUILD_DAY\"/$(BN_0)/g" src/bedrock/RP/manifest.in.json > src/bedrock/RP/manifest.json
	grep version src/bedrock/RP/manifest.json

icon:
	cp expanded_items.png src/java/RP/pack.png
	cp expanded_items.png src/java/DP/pack.png
	cp expanded_items.png src/bedrock/RP/pack_icon.png

package: icon build-number validate
	go run cmd/creeper/main.go -n $(NAME)
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
