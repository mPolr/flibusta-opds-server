package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
)

func unzip(source, dest string) error {
	read, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer read.Close()
	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		name := path.Join(dest, file.Name)
		os.MkdirAll(path.Dir(name), os.ModeDir)
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open")
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()
	return text
}

func findInpFiles(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

func main() {
	fmt.Println("Flibusta OPDS server by mPolr")

	inpxExist := fileExists("./flibusta_fb2_local.inpx")
	fmt.Printf("Checking for 'flibusta_fb2_local.inpx' exist: %v\n", inpxExist)

	if inpxExist {
		unzip("flibusta_fb2_local.inpx", "./inpx")
		fmt.Printf("Checking for 'flibusta_fb2_local.inpx' exist: %v\n", inpxExist)
		version := readFile("./inpx/version.info")
		for _, line := range version {
			fmt.Printf("Found flibusta inpx version %v\n", line)
		}
		for _, s := range findInpFiles("./inpx", ".inp") {
			println(s)
		}
	}
}
