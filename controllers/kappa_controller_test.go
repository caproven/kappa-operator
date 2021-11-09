package controllers

import (
	"testing"
)

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
