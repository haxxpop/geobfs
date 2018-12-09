package main

import (
  "flag"
  "fmt"
  "net"
  "os"
  "strconv"
  geobfs ".."
)

func main() {
  // List of available options.
  port := flag.Int("p", 8080, "The port for the server to listen")
  filename := flag.String("o", "", "The output filename. If not specified, " +
                                   "output will be put in standard output")
  flag.Parse()

  listener, err := net.Listen("tcp", ":" + strconv.Itoa(*port))
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
  }
  defer listener.Close()

  // Listen for the incoming connection.
  conn, err := listener.Accept()
  if err != nil {
    fmt.Fprintln(os.Stderr, err)
    return
  }
  defer conn.Close()

  // Check if the user specify the filename. If so, use the filename
  // instead of the standard output.
  output := os.Stdout
  if *filename != "" {
    output, err = os.OpenFile(*filename, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
      fmt.Fprintln(os.Stderr, err)
      return
    }
  }

  // Wire the connection stream to the output stream.
  geobfs.Deobfuscate(output, conn)
}
