package geobfs

import (
  "bufio"
  "crypto/rand"
  "fmt"
  "io"
  "math/big"
  "strconv"
  "strings"
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
  // The size of sample space used in longitude randomness in the obfuscation.
  LONGITUDE_GRANULARITY = 100000000
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

// Read all lines from src, deobfuscate them, and write them to dst.
func Deobfuscate(dst io.Writer, src io.Reader) error {
  bufreader := bufio.NewReader(src)
  for {
    // Read each line of the stream.
    line, er := bufreader.ReadString('\n')
    if er != nil && er != io.EOF {
      return er
    }
    output, err := deobfuscateLine(line)
    if err != nil {
      return err
    }

    // Write the deobfuscated bytes to dst.
    nw, ew := dst.Write([]byte{ output })
    if ew != nil {
      return ew
    }
    // If this happens, the OS somehow cannot write all the bytes but fail
    // to return an explicit error. We should return an error instead of
    // retrying to write to the destination stream because that can cause
    // an infinite loop.
    if nw != 1 {
      return io.ErrShortWrite
    }

    // If we already reach the end of the stream, return.
    if er == io.EOF {
      return nil
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

// Return a byte slice deobfuscated from the line.
func deobfuscateLine(line string) (byte, error) {
  // Trim a line.
  line = strings.Trim(line, " \n")
  // Trim a prefix.
  coordinate := strings.TrimPrefix(line, "geo:")
  // Split for a latitude and longitude.
  slice := strings.Split(coordinate, ",")
  if len(slice) != 2 {
    return 0, fmt.Errorf("a line %s is malformed", line)
  }

  // Parse for latitude.
  latitude, err := strconv.ParseFloat(slice[0], 64)
  if err != nil {
    return 0, err
  }
  if latitude < LATITUDE_MIN || latitude > LATITUDE_MAX {
    return 0, fmt.Errorf("a latitude %f is not in a valid range", latitude)
  }

  // Parse for longitude.
  // Even if we will discard the longitude, it is better to check that
  // the longitude is valid.
  longitude, err := strconv.ParseFloat(slice[1], 64)
  if err != nil {
    return 0, err
  }
  if longitude < LONGITUDE_MIN || longitude > LONGITUDE_MAX {
    return 0, fmt.Errorf("a longitude %f is not in a valid range", longitude)
  }

  // Deobfuscate a byte from a latitude.
  rng := LATITUDE_MAX - LATITUDE_MIN
  step := float64(rng) / 256
  // Casting float64 to byte will remove noise from the latitude.
  return byte((latitude - float64(LATITUDE_MIN)) / step), nil
}
