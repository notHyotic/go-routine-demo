package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

func tarFiles(files []string, tarName string, progress chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(progress)

	// Create the tar archive
	tarOut, err := os.Create(tarName)
	if err != nil {
		fmt.Println("Error creating tar file:", err)
		return
	}
	defer tarOut.Close()

	tw := tar.NewWriter(tarOut)
	defer tw.Close()

	total := len(files)
	for i, fileName := range files {
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Printf("Error opening %s: %v\n", fileName, err)
			continue
		}

		info, err := file.Stat()
		if err != nil {
			fmt.Printf("Error stating %s: %v\n", fileName, err)
			file.Close()
			continue
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			fmt.Printf("Error creating header for %s: %v\n", fileName, err)
			file.Close()
			continue
		}
		header.Name = filepath.Base(fileName)

		if err := tw.WriteHeader(header); err != nil {
			fmt.Printf("Error writing header for %s: %v\n", fileName, err)
			file.Close()
			continue
		}

		if _, err := io.Copy(tw, file); err != nil {
			fmt.Printf("Error copying %s: %v\n", fileName, err)
			file.Close()
			continue
		}
		file.Close()

		progress <- (i + 1) * 100 / total // send progress
	}
}

func main() {
	files := []string{"file1.txt", "file2.txt", "file3.txt"} // Replace with your real files
	progress := make(chan int)
	var wg sync.WaitGroup
	wg.Add(1)

	go tarFiles(files, "output.tar", progress, &wg)

	// Main thread: Show progress
	for p := range progress {
		fmt.Printf("Progress: %d%%\n", p)
	}

	wg.Wait()
	fmt.Println("All files archived successfully.")
}
