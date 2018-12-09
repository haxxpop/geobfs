package main

import (
  "bufio"
  "encoding/binary"
  "fmt"
  "io"
  "net"
  "os"
)

const (
  SOCKS4_VERSION_NUMBER = 4
  SOCKS4_ESTABLISH_COMMAND = 1
  SOCKS4_REQUEST_HEADER_LEN = 8
)

type SOCKS4Conn struct {
  // The connection connected to the client.
  clientConn net.Conn
  // The connection connected to the server.
  serverConn net.Conn
}

// Handle client connection request. Return a corresponding
// SOCKS4Conn, if success.
func NewSOCKS4Conn(clientConn net.Conn) (SOCKS4Conn, error) {
  // Parse for a remote address.
  address, err := parseRemoteAddr(clientConn)
  if err != nil {
    clientConn.Close()
    return SOCKS4Conn{}, err
  }

  // Since we already have a remote address, we can now create a TCP
  // connection to the remote host.
  serverConn, err := net.Dial("tcp", address)
  if err != nil {
    clientConn.Close()
    return SOCKS4Conn{}, err
  }

  // Send a success response back to the client.
  var response [8]byte
  // Set a request granted status.
  response[1] = 0x5a
  clientConn.Write(response[:])

  conn := SOCKS4Conn{ clientConn, serverConn }
  return conn, nil
}

// Close both server and client connection.
func (conn *SOCKS4Conn) Close() error {
  err := conn.clientConn.Close()
  if err != nil {
    return err
  }

  err = conn.serverConn.Close()
  if err != nil {
    return err
  }

  return nil
}

type SOCKS4Listener struct {
  // The TCP listener.
  netListener net.Listener
}

// Listen for a new connection.
func NewSOCKS4Listener(address string) (SOCKS4Listener, error) {
  netListener, err := net.Listen("tcp", address)

  listener := SOCKS4Listener{ netListener }
  return listener, err
}

// Accept a new connection.
func (listener *SOCKS4Listener) Accept() (SOCKS4Conn, error) {
  // We need to loop until we have a valid connection or there is an error
  // from netListener.
  for {
    clientConn, err := listener.netListener.Accept()
    if err != nil {
      return SOCKS4Conn{}, err
    }

    // Create SOCKS4Conn.
    conn, err := NewSOCKS4Conn(clientConn)
    if err != nil {
      fmt.Fprintln(os.Stderr, err)
      continue
    }
    // Success.
    return conn, nil
  }
}

// Close the listener.
func (listener *SOCKS4Listener) Close() error {
  return listener.netListener.Close()
}

// Return the network address of the listener.
func (listener *SOCKS4Listener) Addr() net.Addr {
  return listener.netListener.Addr()
}

// -------------------- PRIVATE FUNCTIONS --------------------

// Try to parse the remote IP address and TCP port from a reader stream.
func parseRemoteAddr(reader io.Reader) (string, error) {
  // Try to read Fields 1-4 (header) of SOCKS4 request to
  // 1. Verify that the SOCKS version and command code are correct.
  // 2. Find an IP address and port of the remote server.
  var header [SOCKS4_REQUEST_HEADER_LEN]byte
  slice := header[:]

  for len(slice) > 0 {
    n, err := reader.Read(slice)
    slice = slice[n:]
    // Check if we have already read the whole header from the client stream.
    // If so, we can now inspect the header.
    if len(slice) == 0 && err == nil {
      break
    }
    // If there is an error while we have not finished reading the header,
    // close the client connection and return.
    if err != nil {
      return "", fmt.Errorf("socks request is malformed")
    }
  }

  // Check version number
  if version := header[0]; version != SOCKS4_VERSION_NUMBER {
    return "", fmt.Errorf("socks version number %d is not supported", version)
  }
  // Check socks command
  if command := header[1]; command != SOCKS4_ESTABLISH_COMMAND {
    return "", fmt.Errorf("socks command %d is not supported", command)
  }

  // There is another variable-length field called user ID. We need to read
  // and discard it.
  bufreader := bufio.NewReader(reader)
  _, err := bufreader.ReadString(0)
  if err != nil {
    return "", fmt.Errorf("socks request is malformed")
  }

  port := binary.BigEndian.Uint16(header[2:4])
  address := fmt.Sprintf("%d.%d.%d.%d:%d", header[4], header[5],
                                           header[6], header[7], port)
  return address, nil
}
