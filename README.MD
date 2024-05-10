Concurrent File Searcher with Context

build a program that searches for files with a specific name or content in a given directory and its subdirectories concurrently. You'll use Goroutines for concurrent execution, channels for communication between Goroutines, a buffer to limit the number of concurrent file searchers, select to handle timeouts or cancellations, and a mutex to protect shared data structures.

outline of how the project can be structured:

Input Handling: The program should accept input from the user, such as the directory to search, the file name or content to search for, and any search options (e.g., case sensitivity, file extensions to include/exclude).

Concurrent File Search: Use Goroutines to perform file searches concurrently. Each Goroutine should traverse a directory and its subdirectories, searching for files that match the specified criteria (e.g., file name or content). Use channels to send search results (e.g., file paths) back to the main Goroutine.
Buffered Channels: Implement a buffer to limit the number of concurrent file searchers. This helps control resource usage and prevents overwhelming the system with too many concurrent searches.

Context: Use the context package to manage the lifecycle of the file search operation. Pass a context to each Goroutine and use it to handle timeouts or cancellations gracefully. For example, if the user decides to cancel the search or if a timeout occurs, the program should stop searching for files and return the results found so far.
Select: Use a select statement to wait for results from multiple Goroutines concurrently. This allows the program to handle multiple search operations simultaneously and efficiently aggregate results as they become available.

Mutex: Use a mutex to protect shared data structures, such as a map or slice used to store search results. Since multiple Goroutines will be accessing and updating the shared data concurrently, it's essential to use mutexes to prevent data races and ensure thread safety.

Output Display: Once all file searches are complete, display the search results to the user in a clear and organized manner. This could include printing the paths of matching files, along with any additional information (e.g., file size, last modified timestamp).

hands-on experience with concurrent programming concepts in Go, including Goroutines, channels, buffer, select, mutex, and context. You'll also gain practical experience in building concurrent applications that can efficiently handle concurrent tasks and gracefully handle timeouts or cancellations.
