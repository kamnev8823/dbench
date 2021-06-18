package terminal

import "fmt"

type Terminal struct {
	History []string
	Cursor  string
}

func (t *Terminal) PrintCursor() {
	fmt.Print(t.Cursor)
}

func (t *Terminal) SaveHistory(s string) {
	t.History = append(t.History, s)
}

func (t *Terminal) PrintHistory() {
	fmt.Print("\n\tHistory:")
	for k, v := range t.History {
		fmt.Printf("\t%v. %v\n", k, v)
	}
}
