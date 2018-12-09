CLIENT_FILES=$(shell find client -type f)
SERVER_FILES=$(shell find server -type f)
COMMON_FILES=geobfs.go

all: geobfs-client geobfs-server

geobfs-client: $(CLIENT_FILES) $(COMMON_FILES)
	cd client && go build -o ../geobfs-client

geobfs-server: $(SERVER_FILES) $(COMMON_FILES)
	cd server && go build -o ../geobfs-server

clean:
	rm geobfs-client geobfs-server

.PHONY: all clean
