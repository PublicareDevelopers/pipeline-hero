package slack

import "testing"

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
