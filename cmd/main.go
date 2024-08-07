package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	pietvm "github.com/mattellis91/piet-vm"

	tea "github.com/charmbracelet/bubbletea"
)

var verbose = flag.Bool("v", false, "verbose")

func main() {
	flag.Parse()

	fmt.Printf("%v", flag.Args())

	reader, err := os.Open("../testimages/hello_world.png")
	if err != nil {
		log.Fatal(err)
	}

	im, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%v", im)

	p := tea.NewProgram(newModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running tui: %v", err)
		os.Exit(1)
	}


	in := pietvm.New(im)
	in.Run()
	println()
	
	// if len(flag.Args()) != 1 {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

}