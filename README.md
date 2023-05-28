# Alphabet-Counter
A simple gRPC server that maintains the count of alphabets it has received over time.
Client sends 4096 requests/sec to the server while server renders the max frequency alphabet,
its frequency along with the total number of messages it has received since the start of the server.
Here is how the server output looks like...\
![img](https://i.ibb.co/27PC054/2023-05-28-23-14-48-Alphabet-Count-stats-renderer-go.png)\
For the client side we have a little bit of different output
where we simply display the completed requests along with the time it took to complete those requests.
Here is how it looks like...\
![img](https://i.ibb.co/J7WR2sH/2023-05-28-23-15-12-Alphabet-Count-stats-renderer-go.png)
## Prerequisite
The project was created with `go1.20.4`, `libprotoc 23.2` and lastly `makefile(for easy builds and protobufs compilation)`. Make sure you have all of these installed on your machine before proceeding. 
**Important**, the project was created and tested on `Windows 11`, although, make file will run across different platforms keep an eye on errors that might pop up in case make file is not able to compile binaries properly.\
## Running Project
Once you have installed everything required. Simply type `make help` which will display all the commands configure with makefile...\
![img](https://i.ibb.co/Y7Mx3Q5/asd.png)
All of these commands are self-explanatory, simply run `make test` to run server tests, afterward run `make all` which will create server and client binaries and generate proto files for the written contracts. Binaries are places inside `bin` folder of the project...\
![img](https://i.ibb.co/bQ3Pfwj/asd.png)\
You can simply change into the server directory and run it `allow` on Windows firewall defender will pop up simply whitelist the server through it. Next start client, and it will start sending requests to the server. You can start multiple clients if you want.
## General design philosophy
### Data Structure for efficient writes and read
Before getting started with the implementation,
my initial focus was on the core problem of rating words based on their frequency.
The first data structure that came to my mind was `max heap` since insertion will cost `O(log n)` but `peek` itself will be very cheap`O(1)`
but an issue with heap was updates if we need to update an existing value inside heap it will get expensive in worse case like if the element we want to update is at the end of the heap or middle.
Moving on,
I decided
to further simplify the solution
since we know for sure
that Alphabets provided
as input will be 26 how about we create a slice of structs
and simply merge sort them before rendering the alphabet count.
`O(nlogn)` of merge sort won't be a big deal since in the end list will be limited to 26 in size as per the count of alphabets.
Before moving on with its implementation, I thought of an even simpler implementation
focused around buckets.
So what if we create an array of fixed size i.e.,
26
where each index represents an alphabet
offset from its ASC11 representation by a difference of -65 i.e. `A` whose ASCII is `65` if an offset of `-65` is applied it becomes
`0` and can easily be mapped to a linear array of 26 indexes.
Each index mapped to an alphabet while the number value in that index represents the frequency of that alphabet.
We have a global struct that keeps track of the highest frequency alphabet and its frequency.
On each writing we check if the new written alphabet has higher 
frequency than the current if yes simply replace the new alphabet as the one with the highest frequency.
With this we have a solution that's `O(1)` in insertion, update and read.
As far
as space complexity is concerned we can say its `O(1)` since the alphabet count caps the input set to 26 alphabets.
### Data Consistency Problems/Solution
Another problem in the back of my head was data consistency related issues.
Since my plan was to use alphabet as a global counter accessible via each goroutine spawned per gRPC request,
it was for sure going to result in data consistency issues.
My initial plan was to use mutexes/locks for ensuring we control multi-writes and keep data consistent
but after going through diving into `Go` design philosophies and its focus around `channels` I decided
to add a channel dedicated for `AlphabetCounter` which resolved my consistency related issues.
To keep the channel unblocked, I buffered it with a value of `4096`string types.
### Random Generation
Another problem faced while working was Random numbers generation.
While running inside routine, initially I was seeding it multiple times,
which were leading to similar numbers being generated every time.
To find a way around this problem, I created a singleton of random which was seeded on first request.
To protect singleton across `goroutines` I added a lock for safe initialization.
### Client Architecture Scale
Right now client architecture is pretty simple,
after each second we spawn a new `goroutine` which sends 4096 requests and verifies it.
What if we plan to further increase the client request count.
I was thinking of something that can be effectively used to stress test the gRPC capabilities.
For this I had multiple approaches in mind,
but the simplest one was based on controller/worker architecture using `OS pipes`for inter process communication.
Here is how it may look diagrammatically in action...\
![img](https://i.ibb.co/FBGjJqg/1.png)\
So the concept is basic if we have request 8192, we cap each client with 4096 requests, which means the primary
client splits it into N(Should not increase processor core count) number of clients provided via command line arguments.
So each worker further splits that into 4096/(CPU-CORE*2)
routines thus increasing the overall throughput of the application.
The number for each client i.e. 4096 is hypothetical
since right now each go routine can easily handle it under 10 nanoseconds in reality,
4096 should be a pretty big number for this approach to be adopted practically, for now simpler approach implemented here will work fine.