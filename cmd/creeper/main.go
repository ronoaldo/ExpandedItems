package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	baseName = flag.String("n", "custom", "The package base name.")
	icon     = flag.String("i", "icon.png", "The package base icon to be used on all packs.")

	datapackDir = flag.String("d", filepath.Join("output", "datapacks"),
		"The directory to write datapacks to.")
	resourcePackDir = flag.String("r", filepath.Join("output", "resourcepacks"),
		"The directory to write resourcepacks to.")
	bedrockPackDir = flag.String("b", filepath.Join("output", "mcpacks"),
		"The directory to write Bedrock resourcepacks to.")
)

var (
	buildDay    string
	buildNumber string
)

func init() {
	buildDay, buildNumber = calculateBuildNumber()
}

func main() {
	flag.Parse()

	log.Printf("Building packs (buildDay=%s; buildNumber=%s)", buildDay, buildNumber)

	createZipFile("datapack", "src/java/DP", *datapackDir, *baseName+"-datapack.zip")
	createZipFile("resourcepack", "src/java/RP", *resourcePackDir, *baseName+"-resourcepack.zip")
	createZipFile("addons", "src/bedrock/RP", *bedrockPackDir, *baseName+"-resourcepack.mcpack")
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

	zipFile := zip.NewWriter(fd)

	// Add all files from srcDir
	log.Printf("Adding files from %s folder ...", srcDir)
	srcFs := os.DirFS(srcDir)
	if err := fs.WalkDir(srcFs, ".", addFilesToZip(srcFs, zipFile)); err != nil {
		log.Fatalf("Error creating %v: %v", packType, err)
	}
	// Add icon files to pack
	addIconToZip(packType, zipFile)

	zipFile.Close()
}

func addIconToZip(packType string, zipFile *zip.Writer) {
	log.Printf("Adding icon ...")
	iconHeader := zip.FileHeader{
		Name:   "pack.png",
		Method: zip.Deflate,
	}
	if packType == "addons" {
		iconHeader.Name = "pack_icon.png"
	}
	iconFD, err := os.Open(*icon)
	if err != nil {
		log.Fatalf("Error opening icon file '%v' for writting: %v", *icon, err)
	}
	defer iconFD.Close()
	iconW, err := zipFile.CreateHeader(&iconHeader)
	if err != nil {
		log.Fatalf("Error creating icon in zip: %v", err)
	}
	if _, err := io.Copy(iconW, iconFD); err != nil {
		log.Fatalf("Error writing icon to zip: %v", err)
	}
}

func addFilesToZip(srcDir fs.FS, zipFile *zip.Writer) func(path string, d fs.DirEntry, err error) error {
	return func(path string, d fs.DirEntry, err error) error {
		var (
			info fs.FileInfo
			fr   fs.File

			zipHeader *zip.FileHeader
			fw        io.Writer
		)

		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
		// Fetch dir entry info and validate if it can be zipped
		if info, err = d.Info(); err != nil {
			return err
		}
		if !info.IsDir() && !info.Mode().IsRegular() {
			return errors.New("creeper: unsupported file type")
		}
		// Create a zip FileInfoHeader for the path
		if zipHeader, err = zip.FileInfoHeader(info); err != nil {
			return err
		}
		zipHeader.Name = path
		if d.IsDir() {
			zipHeader.Name += "/"
		}
		zipHeader.Method = zip.Deflate
		// Add the header to zipFile
		if fw, err = zipFile.CreateHeader(zipHeader); err != nil {
			return err
		}

		// Nothing to copy, just the dir entry
		if d.IsDir() {
			return nil
		}
		// Read src file from fs
		if fr, err = srcDir.Open(path); err != nil {
			return err
		}
		defer fr.Close()
		b, err := io.ReadAll(fr)
		if err != nil {
			return err
		}
		// Parse to check if is a valid JSON file and add version info if needed
		if strings.HasSuffix(path, ".json") {
			if err = isValidJSON(b); err != nil {
				return fmt.Errorf("creeper: invalid JSON file: %v: %v", path, err)
			}
			b = replaceVersionString(b)
		}
		// Write resulting bytes to zip
		_, err = fw.Write(b)
		return err
	}
}

func replaceVersionString(b []byte) []byte {
	b = bytes.ReplaceAll(b, []byte("\"BUILD_DAY\""), []byte(buildDay))
	b = bytes.ReplaceAll(b, []byte("\"BUILD_NUMBER\""), []byte(buildNumber))
	return b
}

func isValidJSON(b []byte) error {
	j := make(map[string]interface{})
	return json.Unmarshal(b, &j)
}

func calculateBuildNumber() (string, string) {
	now := time.Now()
	buildDay := now.YearDay()
	buildTime := 60*now.Hour() + now.Minute()
	buildSeconds := now.Second() % 10
	return fmt.Sprintf("%d", buildDay), fmt.Sprintf("%d%d", buildTime, buildSeconds)
}

func makeDirs(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		panic(err)
	}
}
