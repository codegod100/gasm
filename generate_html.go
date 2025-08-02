//go:build ignore

package main

import (
	"context"
	"os"
)

func main() {
	// Generate the HTML file using templ
	file, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Render the ChatPage template to the file
	err = ChatPage().Render(context.Background(), file)
	if err != nil {
		panic(err)
	}
}
