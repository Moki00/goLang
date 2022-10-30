package main

import (
    "fmt"
)

func main() {
    fmt.Println(Points([]string{"1:0","2:0","3:0","4:0","2:1","3:1","4:1","3:2","4:2","4:3"}))
}

func Points(games []string) int {
  
  sum:=0
  
  for _, v := range games {

    var s = int(v[0])
    sum+= 1
    fmt.Println(s)
    fmt.Printf("we see %v\n", v)
  }
  
  return sum
}