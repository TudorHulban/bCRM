package main

import "context"

func main() {
	_, cancel := context.WithCancel(context.Background()) // creating context for app
	defer cancel()
}
