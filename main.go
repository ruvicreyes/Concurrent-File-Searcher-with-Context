package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type model struct{}

func (m model) getHome() string {
	//this will determine the path of read
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	//get Main Directory
	main := filepath.Join(home, "go-concurrent")

	return main
}

func (m model) viewer(filesChan <-chan string, ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context canceled, exiting")
			return
		//check if there is still item in channel
		case val, ok := <-filesChan:
			if !ok {
				fmt.Println("No more file found")
				return
			} else {
				fmt.Println(val)
			}
		}
	}


}

func (m model) files(subDir <-chan string, input string) <-chan string {
	out := make(chan string)

	go func() {
		for s := range subDir {
			// Open the subDirectory
			d, err := os.Open(s)
			if err != nil {
				fmt.Println(err)
			}
			defer d.Close()

			// Read the directory contents
			files, err := d.Readdir(-1)
			if err != nil {
				fmt.Println(err)
			}

			// Send each file name through the channel
			for _, file := range files {
				// Check if the file is a regular file
				if file.Mode().IsRegular() {
					if strings.ToUpper(file.Name()) == input {
						out <- "Dir of:" + s + "\n" + "Items Found: " + file.Name()
					}
				}
			}
		}
		close(out)
	}()

	return out
}

func (m model) subDir(mainDir string) <-chan string {
	out := make(chan string)

	go func() {
		// Open the directory
		d, err := os.Open(mainDir)
		if err != nil {
			panic(err)
		}
		defer d.Close()

		// Read the directory contents
		infos, err := d.Readdir(-1)
		if err != nil {
			panic(err)
		}
		for _, sub := range infos {
			//check if the sub is a directory
			if sub.IsDir() {
				out <- filepath.Join(mainDir, sub.Name())
			}
		}
		close(out)
	}()

	return out
}

func (m model) fileSearcher(filename string) {

	// Create a context with cancellation capability
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() //to endure cancel() if main exits

	//get mainDir
	main := m.getHome()

	//get subdirectories
	subsDirChannel := m.subDir(main)

	//Create a channel for list of files
	filesChannel := m.files(subsDirChannel, filename)

	//view Files
	m.viewer(filesChannel, ctx)

}

func (m model) getinput() (filename string, err error) {
	var input string

	fmt.Println("Enter filename:")
	_, err = fmt.Scanln(&input)
	input = strings.ToUpper(input)
	if err != nil {
		return
	}

	return input, err
}

func main() {
	//initialize the model
	m := model{}

	//input Handling
	filename, err := m.getinput()
	if err != nil {
		panic("Error in input")
	}
	// Concurrent File Search
	m.fileSearcher(filename)

}

//review Buffered channels if possible to implement
//Mutex: Use a mutex to protect shared data structures, such as a map or slice used to store search results.


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


