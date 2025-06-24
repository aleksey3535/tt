package main

import "task/internal/app"

func main() {
	app := app.New()
	if err := app.Run(); err != nil {
		panic(err)
	}
}