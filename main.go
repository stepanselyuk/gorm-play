package main

import "github.com/stepanselyuk/gorm-play/cmd"

//go:generate go-bindata -pkg db_mappings -o db_mappings/bindata.go db_mappings/_templates/

func main() {
	cmd.Execute()
}
