package main

import "testing"

func TestMainStartup(t *testing.T) {
    defer func() {
        if r := recover(); r != nil {
            t.Errorf("main panicked: %v", r)
        }
    }()
    go main()
}
