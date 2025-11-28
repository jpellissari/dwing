package main

import "jpellissari/dwing/cmd"

func main() {
	cmd := cmd.NewCmdRoot()

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
