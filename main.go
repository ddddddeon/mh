package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
)

func main() {
	mdPath := flag.String("m", ".", "md path")
	htmlPath := flag.String("h", ".", "html path")
	flag.Parse()
	process(*mdPath, *htmlPath)
}

func process(inPath string, outPath string) {
	s, err := os.Stat(inPath)
	if err != nil {
		fmt.Printf("Could not stat %s\n", inPath)
		return
	}

	if s.IsDir() {
		files, err := os.ReadDir(inPath)
		if err != nil {
			fmt.Printf("Could not read files in directory %s: %s\n", inPath, err)
			return
		}

		if len(files) > 0 {
			for i := range files {
				process(inPath+"/"+files[i].Name(), outPath+"/"+files[i].Name())
			}
		}
	} else {
		mdRegex, _ := regexp.Compile(`.md$`)
		if mdRegex.MatchString(inPath) {
			md, readErr := os.ReadFile(inPath)
			if readErr != nil {
				fmt.Printf("Error opening file %s: %s", inPath, readErr)
			}
			output := markdown.ToHTML(md, nil, nil)

			pathElements := strings.Split(outPath, "/")
			outFile := strings.Replace(pathElements[len(pathElements)-1], ".md", ".html", 1)
			outDir := strings.Join(pathElements[:len(pathElements)-1], "/")
			outPath = outDir + "/" + outFile

			_, dirErr := os.Stat(outDir)
			if dirErr != nil {
				mkDirErr := os.MkdirAll(outDir, os.ModePerm)
				if mkDirErr != nil {
					fmt.Printf("Could not create directory %s: %s\n", outDir, dirErr)
				}
			}

			writeErr := os.WriteFile(outDir+"/"+outFile, output, os.FileMode(int(0770)))
			if writeErr != nil {
				fmt.Printf("Could not write to file %s: %s\n", outPath, writeErr)
			}
			fmt.Printf("%s -> %s\n", inPath, outPath)
		}
	}
}
