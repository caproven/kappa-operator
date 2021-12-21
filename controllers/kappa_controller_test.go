package controllers

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KappaController", func() {
	// kappa := &caproveninfov1.Kappa{}

	Describe("idk", func() {
		It("should say hello", func() {
			Expect(SayHello("Kappa")).To(Equal("Hello, Kappa"))
		})
	})
})

func TestKappaController(t *testing.T) {
}

func TestSayHello(t *testing.T) {
	tests := map[string]string{
		"Kappa":     "Hello, Kappa",
		"Mr. Kappa": "Hello, Mr. Kappa",
		"":          "Hello!",
	}

	for name, expected := range tests {
		actual := SayHello(name)
		if actual != expected {
			t.Errorf("Expected %s, got %s", expected, actual)
		}
	}
}
