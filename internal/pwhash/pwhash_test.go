package pwhash

import "testing"

func TestPasswordHasher(t *testing.T) {
	if _, err := New(1, 1); err == nil {
		t.Fatal("created hasher with insecure parameters")
	}

	ph, err := New(16, 100000)
	if err != nil {
		t.Fatal(err)
	}
	hash, err := ph.HashPassword("foo")
	if err != nil {
		t.Fatal(err)
	}
	if len(hash) < 129 {
		t.Errorf("hash is too short: %s", hash)
	}

	if err := ph.Validate("boo", hash); err == nil {
		t.Errorf("boo was accepted instead of foo")
	}
	if err := ph.Validate("foo", hash); err != nil {
		t.Errorf("foo wasn't accepted: %v", err)
	}

	adminHash := "PBKDF2$sha512$1$l8zGKtxRESq3PA2kFhHRWA==$H3lGMxOt55wjwoc+myeOoABofJY9oDpldJa7fhqdjbh700V6FLPML75UmBOt9J5VFNjAL1AvqCozA1HJM0QVGA=="
	if err := ph.Validate("admin", adminHash); err != nil {
		t.Errorf("admin hash couldn't be validated: %v", err)
	}
}
