package pkg

import "testing"

func TestHello(t *testing.T) {
	got := Hello()
	want := "Hello from pkg!"
	if got != want {
		t.Errorf("Hello() = %q; want %q", got, want)
	}
}
