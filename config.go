package main

import "os"

// should change to the bot
const (
	GH_USERNAME = "2nthony"
	GH_EMAIL    = "hi@2nthony.com"
)

var GH_TOKEN = os.Getenv("GH_TOKEN")

const TARGET_ORG_URL = "https://github.com/wails-template"
