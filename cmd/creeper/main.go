package main

import (
	"archive/zip"
	"flag"
	"log"
	"os"
	"path/filepath"
)

var (
	baseName = flag.String("n", "custom", "The package base name.")

	datapackDir = flag.String("d", filepath.Join("output", "datapacks"),
		"The directory to write datapacks to.")
	resourcePackDir = flag.String("r", filepath.Join("output", "resourcepacks"),
		"The directory to write resourcepacks to.")
	bedrockPackDir = flag.String("b", filepath.Join("output", "mcaddon"),
		"The directory to write Bedrock resourcepacks to.")
)

func main() {
	flag.Parse()

	log.Printf("Creating packages with creeper...")

	createZipFile("datapack", "src/java/DP", *datapackDir, *baseName+"-datapack.zip")
	createZipFile("resourcepack", "src/java/RP", *resourcePackDir, *baseName+"-resourcepack.zip")
	createZipFile("addons", "src/bedrock/RP", *bedrockPackDir, *baseName+"-resourcepack.mcpack")
}

func makeDirs(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}
}

func createZipFile(packType, srcDir, dstDir, fileName string) {
	path := filepath.Join(dstDir, fileName)
	log.Printf("Creating a %s at %v", packType, path)

	// Create destDir if not exists
	makeDirs(dstDir)
	// Open zip file
	fd, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error opening file '%v' for writting: %v", path, err)
	}
	defer fd.Close()
	// Add all files from srcDir
	log.Printf("Adding files from %s folder ...", srcDir)
	zip := zip.NewWriter(fd)
	if err := zip.AddFS(os.DirFS(srcDir)); err != nil {
		log.Fatalf("Error adding files to datapack: %v", err)
	}
	zip.Close()
}
