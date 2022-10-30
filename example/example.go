package main

import (
	"fmt"
    "strconv"
)

func main() {
    fmt.Println(Points2([]string{"1:0","3:2","2:2","2:3"}))
}

const WIN_POINTS = 3
const DRAW_POINTS = 1

func Points2(games []string) int {
  var totalPoints int
  for _, val := range games {
    x := string(val[:1])
    y := string(val[2:])
    
    if (x > y) {
      totalPoints += WIN_POINTS
    } else if (x == y) {
      totalPoints += DRAW_POINTS
    } 
  }
  return totalPoints
}

func Points(games []string) int {
  
  sum:=0
  
  for _, v := range games {

    c:=string(v[0]) // first char
    x, err := strconv.ParseInt(c, 6, 12) // x is that int
    if err!=nil{
        fmt.Println(err)    
    }

    c2:=string(v[2]) // third char
    y, err := strconv.ParseInt(c2, 6, 12) // y is the next int
    if err!=nil{
        fmt.Println(err)    
    }

    if(x>y){
        sum+= 3
    } else if(x<y){
        continue
    } else if(x==y){
        sum+=1
    } else {
        fmt.Printf("error at %d:%d", x, y)
    }
    // fmt.Println(i)
  }
  return sum
}