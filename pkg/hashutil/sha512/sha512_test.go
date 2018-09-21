package sha512

import (
	"encoding/base64"
	"testing"
)

var commonTests = []struct {
	in       string
	expected string
}{
	{"", "z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg=="},
	{"angryMonkey", "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="},
	{"superSecret", "pQaPqt7aC/CThmNsO8xnV+nkLfyJoyqpFzGzmvLivIpjmQXnvJqIULCUOpE+H1f3+p9laadfIkvAxMYZTAxnyQ=="},
	{"never gonna give you up", "ll0dEoj9qD0ALPTTXdjQL+XD4fvmpdAw67Gql5jfTx6cv16SeveRSmRX5a9qmC0GN99PV+1h6jlrwS+NoUq3Ug=="},
}

func TestWriteString(t *testing.T) {
	for _, test := range commonTests {
		hasher := NewStringWriter()

		_, err := hasher.WriteString(test.in)
		if err != nil {
			t.Fatal(err)
		}

		actual := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
		if actual != test.expected {
			t.Errorf("Expected: %s, Actual: %s", test.expected, actual)
		}
	}
}
