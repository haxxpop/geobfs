package geobfs

import (
  "crypto/rand"
  "fmt"
  "io"
  "math/big"
)

const (
  // The default buffer size.
  DEFAULT_BUFFER_LEN = 128
  // Min and max of latitude.
  LATITUDE_MIN = -90
  LATITUDE_MAX = 90
  // Min and max of longitude.
  LONGITUDE_MIN = -180
  LONGITUDE_MAX = 180
  // The size of sample space used in noise randomness in the obfuscation.
  NOISE_GRANULARITY = 10000
)

// Read all bytes from src, obfuscate them, and write them to dst.
func Obfuscate(dst io.Writer, src io.Reader) error {
  // Allocate a buffer to store the read bytes.
  buf := make([]byte, DEFAULT_BUFFER_LEN)
  for {
    nr, er := src.Read(buf)
    if er != nil && er != io.EOF {
      return er
    }
    if nr > 0 {
      // We will obfuscate one byte at a time.
      for _, b := range buf[:nr] {
        line := obfuscateLine(b)
        // Write the obfuscated string to dst.
        nw, ew := dst.Write([]byte(line))
        if ew != nil {
          return ew
        }
        // If this happens, the OS somehow cannot write all the bytes but fail
        // to return an explicit error. We should return an error instead of
        // retrying to write to the destination stream because that can cause
        // an infinite loop.
        if len(line) != nw {
          return io.ErrShortWrite
        }
      }
    }

    // If we already reach the end of the stream, return.
    if er == io.EOF {
      return nil
    }
  }
}

func Deobfuscate(dst io.Writer, src io.Reader) error {
  // Allocate a buffer to store the read bytes.
  read_buf := make([]byte, DEFAULT_BUFFER_LEN)
  for {
    nr, er := src.Read(read_buf)
    if er != nil  {
      // If it is just an EOF, no need to return error.
      if er == io.EOF {
        return nil
      }
      return er
    }
  }
}

// -------------------- PRIVATE FUNCTIONS --------------------

// Return a NL-terminated string representing the obfuscated byte.
func obfuscateLine(val byte) string {
  rng := LATITUDE_MAX - LATITUDE_MIN
  step := float64(rng) / 256
  // Generate noise to be added to the latitude.
  // We know that there is no error because GRANULARITY is greater
  // than zero.
  random, _ := rand.Int(rand.Reader, big.NewInt(NOISE_GRANULARITY))
  noise := float64(random.Int64()) / NOISE_GRANULARITY * step
  // Calculate latitude.
  latitude := float64(val) * step + float64(LATITUDE_MIN) + noise

  // For the longitude, we will just random it uniformly.
  rng = LONGITUDE_MAX - LONGITUDE_MIN
  random, _ = rand.Int(rand.Reader, big.NewInt(LONGITUDE_GRANULARITY))
  longitude := float64(random.Int64()) / LONGITUDE_GRANULARITY *
                                                  float64(rng) + LONGITUDE_MIN
  return fmt.Sprintf("geo:%f,%f\n", latitude, longitude)
}
