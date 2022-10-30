package main

import (
	"fmt"
)

func main() {
    fmt.Println(Rps("rock","scissors"))
    fmt.Println(Rps("scissors","rock"))
    fmt.Println(Rps("rock","rock"))
}

var m = map[string]string{"rock": "paper", "paper": "scissors", "scissors": "rock"}

func Rps(a, b string) string {
  if a == b {
    return "Draw!"
  }
  if m[a] == b {
    return "Player 2 won!"
  }
  return "Player 1 won!"
}

func Rps2(p1, p2 string) string {
	const p string = "paper"
	const r string = "rock"
	const s string = "scissors"
	const p11 string = "Player 1 won!"
	const p21 string = "Player 2 won!"
	const d string = "Draw!"
	var a string = "answer"
  
	switch {
	  case p1==p:
		if p2==p {
			a=d
		} else if p2==r {
			a=p11
		} else {
			a=p21
		}
	  case p1==r:
		if p2==p {
			a=p21
		} else if p2==r {
			a=d
		} else {
			a=p11
		}
	  case p1==s:
		if p2==p {
			a=p11
		} else if p2==r {
			a=p21
		} else {
			a=d
		}
	}
	return a
}