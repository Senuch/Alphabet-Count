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
4096 should be a pretty big number for this approach to be adopted practically.
For now, the simpler approach implemented here will work fine.
## Backend Scaling
If we take a look at the current architecture and it being scaled to a hundred
of nodes, there are 2 major parts of the server that can cause problems which are as follows...
1. gRPC sticky session
2. Distributed Counter(Happens as a result of solving the first issue)
3. Count Display Service
### gRPC Sticky Sessions
gRPC sticky sessions use the same connection through most of the life cycle
of the client requests, and this actually becomes even trickier when the connection is streaming.
Now as the number of clients increases and on top,
processing millions of requests, it will become difficult to manage all these sessions and requests on a single server.
A clear solution is to use multiple gRPC servers with client or discovery focused load balancing.
We can use external load balancer which is continuously fed with performance metrics of each server, and when a client
wants to connect with the server, server gives the client a properly sorted list based on servers metric.
Here is how this looks like in action...\
![img](https://i.ibb.co/fpwyqP0/2.png)\
So External LB is keeps track of the server metrics while also being consumed by the clients for fetching
a suitable servers list for connection.
With this, our clients can easily reach out to the server and get a connection, but
it breaks the counter logic which leads us to the second problem concerning Distributed Counter.
### Distributed Counters
Previously, our counter state was maintained on a single server which is rendered bootless thanks to the server fleet.
Now there are two possible solutions which are as follows...
1. Implementing custom distributed counters
2. Using pre-existing message queues

Also since we don't plan to persist information, and the data itself is limited to 26 alphabets.
Both of these
 make counter implementation a little less tricky.
We can add make use of Kafka.
Since we are considering 100 servers for our scenario,
we can add 26 Kafka topics with multiple partitions per topic.
Each topic will be used by a single alphabet, and multiple writes on a single topic will be entertained
by multiple consumers consuming multiple partitions of a single alphabet topic.
Here is how it looks like...\
![img](https://i.ibb.co/cc5215R/zdsas.png)\
So multiple consumers will be consuming a single topic partitions.
Each consumer will write the incremented count
to a single shared location which can be a cache like redis maintaining count of all 26 Alphabets.
### Count Display Service
This will be the last piece that needs to be refactored.
This service will display the current count of messages along
with the alphabet with the highest frequency.
Because of the scale, our data won't be 
highly consistent but eventually consistent.
If we go with high consistency, it might induce latency which we don't want, so we
will compromise on high consistency and stick with eventually consistent cache and count rendering service.