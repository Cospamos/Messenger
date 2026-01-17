package main

import "test/app"

func main() {
	app.DotEnvInit() // utils.go
	app.RedisInit()  // utils.go

	app.ClearDB() 	 // utils.go

	defer app.StopNgrok() 
	app.UseNgrok()
	// app.GenQR() 
	app.RunServer()
}