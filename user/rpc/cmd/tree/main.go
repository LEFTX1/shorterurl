package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func printTree(path string, level int) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, entry := range entries {
		prefix := strings.Repeat("  ", level)
		fmt.Println(prefix + "|-- " + entry.Name())

		if entry.IsDir() {
			printTree(filepath.Join(path, entry.Name()), level+1)
		}
	}
}

func main() {
	root := "../../" // Replace with the path to your project directory
	fmt.Println(root)
	printTree(root, 0)
}
