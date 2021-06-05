package main

import "fmt"

type Terminal struct {
	history []string
	cursor  string
}

func (t *Terminal) Cursor() {
	fmt.Print(t.cursor)
}

func (t *Terminal) SaveHistory(s string) {
	t.history = append(t.history, s)
}

func (t *Terminal) PrintHistory() {
	fmt.Println("\n\tHistory:")
	for k, v := range t.history {
		fmt.Printf("\t%v. %v\n", k, v)
	}
}
