# Geobfs

## Getting Started
Check out the repository.
```
git clone https://github.com/haxxpop/geobfs.git
cd geobfs
```
Build `geobfs-client` and `geobfs-server`. (You must have already installed Go)
```
make
```
Run the geobfs SOCKS4 proxy server with a specified port number. (Default to 8081, if not specified) 
```
./geobfs-client -p 8081
```
Run the geobfs server with a specified port number. (Default to 8080, if not specified)
```
./geobfs-server -p 8080 -o output.txt # if you want to store the output in a file.
./geobfs-server -p 8080 # if you want to see the output in the standard output.
```
Try sending some message to the geobfs server via the geobfs SOCKS4 proxy server.
```
nc localhost 8080 -X 4 -x localhost:8081
```
Run tests.
```
make check
```
