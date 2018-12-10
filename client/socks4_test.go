package main

import (
  "bytes"
  "errors"
  "testing"
)

func TestParseRemoteAddr(t *testing.T) {
  cases := []struct {
    in []byte
    want_address string
    want_err error
  }{
    { []byte{}, "",
      errors.New("socks request is malformed") },
    { []byte{4}, "",
      errors.New("socks request is malformed") },
    { []byte{4, 1, 0, 0, 0, 0, 0, 0}, "",
      errors.New("socks request is malformed") },
    { []byte{5, 0, 0, 0, 0, 0, 0, 0, 0}, "",
      errors.New("socks version number 5 is not supported") },
    { []byte{3, 0, 0, 0, 0, 0, 0, 0, 0}, "",
      errors.New("socks version number 3 is not supported") },
    { []byte{4, 2, 0, 0, 0, 0, 0, 0, 0}, "",
      errors.New("socks command 2 is not supported") },
    { []byte{4, 1, 0, 0, 0, 0, 0, 0, 0}, "0.0.0.0:0", nil },
    { []byte{4, 1, 0, 80, 127, 0, 0, 0, 0}, "127.0.0.0:80", nil },
  }

  for _, c := range cases {
    got_address, got_err := parseRemoteAddr(bytes.NewReader(c.in))
    if got_address != c.want_address ||
       (got_err == nil && c.want_err != nil) ||
       (got_err != nil && c.want_err == nil) ||
       (got_err != nil && c.want_err != nil &&
        got_err.Error() != c.want_err.Error()) {
      t.Errorf("parseRemoteAddr(%q) == %q, %q (want %q, %q)", c.in,
               got_address, got_err,
               c.want_address, c.want_err)
    }
  }
}
