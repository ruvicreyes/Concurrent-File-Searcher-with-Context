package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type model struct{}

func (m model) getHome() string {
	//this will determine the path of read
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to read directory %s: %v", home, err)
	}
	fmt.Println(home)
	//get Main Directory
	main := filepath.Join(home)

	return main
}

func (m model) viewer(filesChan <-chan string) {
	if len(filesChan) == 0 {
		fmt.Println("No files found")
		return
	} else {
		for val := range filesChan {
			fmt.Println(val)
		}
	}
}

func (m model) files(subDir <-chan string, input string, ctx context.Context) <-chan string {
	out := make(chan string, 10) // Buffered channel to limit concurrent file searches

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, exiting")
			return
		default:
			defer close(out)
			for s := range subDir {
				// Open the subDirectory
				folder, err := os.Open(s)
				if err != nil {
					log.Printf("Failed to open directory %s: %v", s, err)
					return
				}
				defer folder.Close()

				// Read the directory contents
				files, err := folder.Readdir(-1)
				if err != nil {
					log.Printf("Failed to read directory %s: %v", s, err)
					return
				}
				// Send each file name through the channel
				for _, file := range files {
					// Check if the file is a regular file
					if file.Mode().IsRegular() && strings.ToUpper(file.Name()) == input {
						out <- fmt.Sprintf("Dir of: %s\nItems Found: %s", s, file.Name())
					}
				}
			}
		}
	}()

	return out
}

func (m model) subDir(mainDir string, subDirChan chan<- string, wg *sync.WaitGroup) {
	// Read the directory contents
	entries, err := os.ReadDir(mainDir)
	if err != nil {
		fmt.Println(err)
	}

	for _, entry := range entries {
		path := filepath.Join(mainDir, entry.Name())
		if entry.IsDir() {
			subDirChan <- path // Send the subdirectory to the channel
			fmt.Println("Directory:", path)

			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				m.subDir(p, subDirChan, wg)
			}(path)
		}
	}
		// else {
		// 	fmt.Println("File:", path)
		// }
	
	//out <- filepath.Join(mainDir, sub.Name())

}

func (m model) fileSearcher(filename string) {

	// Create a context with cancellation capability
	//ctx, cancel := context.WithCancel(context.Background())

	//get mainDir
	main := m.getHome()

	//subDirChan := make(chan string, 10) // Buffered channel to limit concurrent subdirectory searches
	subDirChan := make(chan string) // Channel to store subdirectories
	var wg sync.WaitGroup

	//get subdirectories
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.subDir(main, subDirChan, &wg)
	}()
	// Start a goroutine to close the channel once all directories are processed
	go func() {
		wg.Wait()
		close(subDirChan)
	}()

	for dir := range subDirChan {
		fmt.Println("Directory:", dir)
	}

	//it didnt go too Go-Concurrent folder

	//Create a channel for list of files
	// filesChannel := m.files(subsDirChannel, filename, ctx)
	// time.Sleep(2 * time.Second)
	// cancel() // after 2  seconds cancel the context

	// //view Files
	// m.viewer(filesChannel)

}

func (m model) getinput() (filename string, err error) {
	var input string
	fmt.Println("Enter filename:")
	_, err = fmt.Scanln(&input)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(input), nil
}

func main() {
	//initialize the model
	m := model{}

	//input Handling
	filename, err := m.getinput()
	if err != nil {
		log.Fatalf("Error in input: %v", err)
	}
	fmt.Println("Searching........................................")
	// Concurrent File Search
	m.fileSearcher(filename)

	fmt.Println("Search completed---------------------------------")

}

//Context in file search: if the user decides to cancel the search or if a timeout occurs, the program should stop searching for files and return the results found so far.
// Feature of All directories in system using Recursive loop

/* Concurrent File Searcher with Context

In this project, you'll build a program that searches for files with a specific name or content in a given directory and its subdirectories concurrently. You'll use Goroutines for concurrent execution, channels for communication between Goroutines, a buffer to limit the number of concurrent file searchers, select to handle timeouts or cancellations, and a mutex to protect shared data structures.

Here's an outline of how the project can be structured:

Input Handling: The program should accept input from the user, such as the directory to search, the file name or content to search for, and any search options (e.g., case sensitivity, file extensions to include/exclude).

Concurrent File Search: Use Goroutines to perform file searches concurrently. Each Goroutine should traverse a directory and its subdirectories, searching for files that match the specified criteria (e.g., file name or content). Use channels to send search results (e.g., file paths) back to the main Goroutine.
Buffered Channels: Implement a buffer to limit the number of concurrent file searchers. This helps control resource usage and prevents overwhelming the system with too many concurrent searches.

Context: Use the context package to manage the lifecycle of the file search operation. Pass a context to each Goroutine and use it to handle timeouts or cancellations gracefully. For example, if the user decides to cancel the search or if a timeout occurs, the program should stop searching for files and return the results found so far.
Select: Use a select statement to wait for results from multiple Goroutines concurrently. This allows the program to handle multiple search operations simultaneously and efficiently aggregate results as they become available.

Mutex: Use a mutex to protect shared data structures, such as a map or slice used to store search results. Since multiple Goroutines will be accessing and updating the shared data concurrently, it's essential to use mutexes to prevent data races and ensure thread safety.

Output Display: Once all file searches are complete, display the search results to the user in a clear and organized manner. This could include printing the paths of matching files, along with any additional information (e.g., file size, last modified timestamp).

By working on this project, you'll get hands-on experience with concurrent programming concepts in Go, including Goroutines, channels, buffer, select, mutex, and context. You'll also gain practical experience in building concurrent applications that can efficiently handle concurrent tasks and gracefully handle timeouts or cancellations.
*/
