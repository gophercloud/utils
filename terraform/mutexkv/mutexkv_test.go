package mutexkv

import (
	"testing"
	"time"
)

func TestMutexKVLock(t *testing.T) {
	s := struct {
		mkv *MutexKV
	}{
		mkv: NewMutexKV(),
	}

	s.mkv.Lock("foo")

	doneCh := make(chan struct{})

	go func() {
		s.mkv.Lock("foo")
		close(doneCh)
	}()

	select {
	case <-doneCh:
		t.Fatal("Second lock was able to be taken. This shouldn't happen.")
	case <-time.After(50 * time.Millisecond):
		// pass
	}
}

func TestMutexKVUnlock(t *testing.T) {
	s := struct {
		mkv *MutexKV
	}{
		mkv: NewMutexKV(),
	}

	s.mkv.Lock("foo")
	s.mkv.Unlock("foo")

	doneCh := make(chan struct{})

	go func() {
		s.mkv.Lock("foo")
		close(doneCh)
	}()

	select {
	case <-doneCh:
		// pass
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Second lock blocked after unlock. This shouldn't happen.")
	}
}

func TestMutexKVDifferentKeys(t *testing.T) {
	s := struct {
		mkv *MutexKV
	}{
		mkv: NewMutexKV(),
	}

	s.mkv.Lock("foo")

	doneCh := make(chan struct{})

	go func() {
		s.mkv.Lock("bar")
		close(doneCh)
	}()

	select {
	case <-doneCh:
		// pass
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Second lock on a different key blocked. This shouldn't happen.")
	}
}
