all: package

package:
	go run cmd/creeper/main.go -n copper-items

test:
	go run cmd/creeper/main.go -n copper-items \
		-d output/data/world/datapacks \
		-r output/resourcepacks \
		-b output/data/plugins/Geyser-Spigot/packs
	cp src/bedrock/geyser_mappings.json output/data/plugins/Geyser-Spigot/custom_mappings

clean:
	rm -rf output