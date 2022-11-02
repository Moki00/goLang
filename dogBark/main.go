package main

import (
	"fmt"
)

func main() {
	poodle := Dog{"Poodle", 10, "Woof!"}
	fmt.Println(poodle) // Poodle 10 Woof
	fmt.Printf("%+v\n", poodle)  // all k&v
	fmt.Printf("Breed: %v\nWeight: %v\n", poodle.Breed, poodle.Weight) // specific

	poodle.Speak()
	poodle.Sound = "Arf!"
	poodle.Speak()
	poodle.SpeakThreeTimes()
}

// Dog is a struct
type Dog struct {
	Breed  string
	Weight int
	Sound  string
}

// Speak is the function. d Dog is the receiver
func (d Dog) Speak() {
	fmt.Println(d.Sound)
}

// SpeakThreeTimes is how the dog speaks loudly
func (d Dog) SpeakThreeTimes() {
	d.Sound = fmt.Sprintf("%v %v %v", d.Sound, d.Sound, d.Sound)
	fmt.Println(d.Sound)
}
