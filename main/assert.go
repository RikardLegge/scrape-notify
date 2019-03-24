package main

import "testing"

func assertFalse(t *testing.T, val interface{}, message ...string) {
	switch val.(type) {
	case error:
		{
			if val == nil {
				t.Fatal(val, message)
			}
		}
	case bool:
		{
			if val == true {
				t.Fatal(val, "Expected false got true", message)
			}
		}
	}
}

func assert(t *testing.T, val interface{}, message ...string) {
	switch val.(type) {
	case error:
		{
			if val != nil {
				t.Fatal(val, message)
			}
		}
	case bool:
		{
			if val == false {
				t.Fatal(val, "Expected true got false", message)
			}
		}
	}
}
