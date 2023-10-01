package notifier

import "testing"

func TestNew(t *testing.T) {
	handler, err := New("slack")
	if err != nil {
		t.Fatal(err)
	}

	if handler.Client == nil {
		t.Fatal("client is nil")
	}

	_, err = New("invalid")
	if err == nil {
		t.Fatal("expected error")
	}
}
