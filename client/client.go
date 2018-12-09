package main

import (
  "flag"
  "fmt"
  "os"
  "strconv"
  geobfs ".."
)

func main() {
  // List of available options.
  port := flag.Int("p", 8081, "The port for the server to listen")
  flag.Parse()

  // Connect to the remote host.
  listener, err := NewSOCKS4Listener(":" + strconv.Itoa(*port))
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
  }
  defer listener.Close()

  for {
    conn, err := listener.Accept()
    if err != nil {
      fmt.Fprintln(os.Stderr, err)
      break
    }
    // We need to use Go routine so that we can handle multiple connections
    // simultaneously.
    go func() {
      // Wire the server connection to the client connection.
      geobfs.Obfuscate(conn.serverConn, conn.clientConn)
      conn.Close()
    }()
  }
}
