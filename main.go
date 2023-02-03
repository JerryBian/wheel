package main

import (
	"embed"
	"log"
	"os"

	//C "github.com/kelseyhightower/envconfig"
)

//go:embed static/style.min.css static/script.min.js static/fonts/*
var f embed.FS
var AppVer = "1.0"
var GitHash = "1234567"
var BuildTime = "2022-12-30"

func main() {
	log.SetOutput(os.Stdout)

	log.Println("Hello World!")
}