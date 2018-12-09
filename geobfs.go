package geobfs

import "io"

func Obfuscate(dst io.Writer, src io.Reader) error {
  _, err := io.Copy(dst, src)
  return err
}

func Deobfuscate(dst io.Writer, src io.Reader) error {
  _, err := io.Copy(dst, src)
  return err
}
