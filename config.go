package main

import "os"

const (
	GH_USERNAME = "2nthony"
	GH_EMAIL    = "hi@2nthony.com"
)

var GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")

const TARGET_ORG_URL = "https://github.com/wails-template"
