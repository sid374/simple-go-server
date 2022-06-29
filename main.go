package main

import (
	"sid/simpleserver/hello/mailmerge"
)

func main() {
	mailmerge.GetTemplates("Hello <FirstName> <LastName>")
}