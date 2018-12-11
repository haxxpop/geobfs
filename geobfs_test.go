package geobfs

import (
  "errors"
  "testing"
)

func TestObfuscateAndDeobfuscateLine(t *testing.T) {
  // Test round-trip.
  for _, want_byte := range []byte{ 0, 10, 128, 255 } {
    line := obfuscateLine(want_byte)
    got_byte, got_err := deobfuscateLine(line)
    if got_err != nil || got_byte != want_byte {
      t.Errorf("deobfuscateLine(obfuscateLine(%q)) == %q, %q (want %q, nil)",
               want_byte, got_byte, got_err, want_byte)
    }
  }

  // Please note that we don't need to test the function obfuscateLine
  // separately because it is a nondeterministic function and we have already
  // tested deobfuscateLine and the round-trip one. That is more than enough
  // to test obfuscateLine. Another reason is that, since it is a
  // nondeterministic function, the only way to test it is to re-implement
  // the function deobfuscateLine here, which doesn't make sense.

  // Test deobfuscation.
  cases := []struct {
    in string
    want_byte byte
    want_err error
  }{
    { "geo:-21.557039,-179.363102", 97, nil },
    { "geo:-20.914523,-164.776363", 98, nil },
    { "abc", 0, errors.New("a line abc is malformed") },
    { "geo:", 0, errors.New("a line geo: is malformed") },
    { "geo:1", 0, errors.New("a line geo:1 is malformed") },
    { "geo:1,2,3", 0, errors.New("a line geo:1,2,3 is malformed") },
    { "geo:-91,0", 0,
      errors.New("a latitude -91.000000 is not in a valid range") },
    { "geo:91,0", 0,
      errors.New("a latitude 91.000000 is not in a valid range") },
    { "geo:0,-181", 0,
      errors.New("a longitude -181.000000 is not in a valid range") },
    { "geo:0,181", 0,
      errors.New("a longitude 181.000000 is not in a valid range") },
  }

  for _, c := range cases {
    got_byte, got_err := deobfuscateLine(c.in)
    if got_byte != c.want_byte ||
       (got_err == nil && c.want_err != nil) ||
       (got_err != nil && c.want_err == nil) ||
       (got_err != nil && c.want_err != nil &&
        got_err.Error() != c.want_err.Error()) {
      t.Errorf("deobfuscateLine(%q) == %q, %q (want %q, %q)", c.in,
               got_byte, got_err,
               c.want_byte, c.want_err)
    }
  }

  // For some cases, we know that there must be an error but we don't know
  // what error message it will produce. So, we list those cases here.
  system_err_cases := []string{
    "geo:a,b",
    "geo:0,b",
    "geo:a,0",
    "geo:1.11.1,2.101.1",
  }

  for _, c := range system_err_cases {
    got_byte, got_err := deobfuscateLine(c)
    if got_byte != 0 || got_err == nil {
      t.Errorf("deobfuscateLine(%q) == %q, %q (want 0, <error>)", c,
               got_byte, got_err)
    }
  }
}
