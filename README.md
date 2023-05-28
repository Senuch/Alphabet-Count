# Alphabet-Counter
A simple gRPC server that maintains the count of alphabets it has received over time. Client send 4096 requests/sec to the server while server renders the max frequency alphabet, its frequency along with the total number of messages it has received since the start of the server.
Here is how the server output looks like...\
![img](https://i.ibb.co/27PC054/2023-05-28-23-14-48-Alphabet-Count-stats-renderer-go.png)\
For the client side we have a little bit of different output where we simply display the completed requests along with the time it took to complete those requests. Here is how it looks like...\
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
