package sha512

import (
	"testing"
)

var commonTests = []struct {
	in       string
	expected string
}{
	{in: "", expected: "z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg=="},
	{in: "angryMonkey", expected: "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="},
	{in: "superSecret", expected: "pQaPqt7aC/CThmNsO8xnV+nkLfyJoyqpFzGzmvLivIpjmQXnvJqIULCUOpE+H1f3+p9laadfIkvAxMYZTAxnyQ=="},
	{in: "never gonna give you up", expected: "ll0dEoj9qD0ALPTTXdjQL+XD4fvmpdAw67Gql5jfTx6cv16SeveRSmRX5a9qmC0GN99PV+1h6jlrwS+NoUq3Ug=="},
}

func TestWriteString(t *testing.T) {
	for _, test := range commonTests {
		hasher := New()

		_, err := hasher.WriteString(test.in)
		if err != nil {
			t.Fatal(err)
		}

		actual := hasher.String()
		if actual != test.expected {
			t.Errorf("Expected: %s, Actual: %s", test.expected, actual)
		}
	}
}

func TestReset(t *testing.T) {
	hasher := New()

	for _, test := range commonTests {
		_, err := hasher.WriteString(test.in)
		if err != nil {
			t.Fatal(err)
		}

		actual := hasher.String()
		hasher.Reset()

		if actual != test.expected {
			t.Errorf("Expected: %s, Actual: %s", test.expected, actual)
		}
	}
}
