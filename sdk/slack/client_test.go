package slack

import "testing"

var testChannel = "CH7EEHVDW"

func Test_NewTestClient(t *testing.T) {
	_, err := NewTestClient()
	if err != nil {
		t.Error(err)
	}
}

func TestSendProgressSlackMessage(t *testing.T) {
	client, err := NewTestClient()
	if err != nil {
		t.Error(err)
	}

	err = client.SendProgressSlackMessage("test")
	if err != nil {
		t.Error(err)
	}
}

func TestList(t *testing.T) {
	client, err := NewTestClient()
	if err != nil {
		t.Error(err)
	}

	channels := client.list()

	if err != nil {
		t.Error(err)
	}

	if len(channels) == 0 {
		t.Error("expected channels, got 0")
	}
}
