package main

import (
	"fmt"
	"strconv"
)

/*
  Given an array of integers as strings and numbers, return the sum of the array values as if all were numbers.
*/

func main() {
  fmt.Println(SumMix([]any{9, 3, "7", "3"})) // 22
  fmt.Println(SumMix([]any{"5", "0", 9, 3, 2, 1, "9", 6, 7}))
  fmt.Println(SumMix([]any{"3", 6, 6, 0, "5", 8, 5, "6", 2,"0"}))
  fmt.Println(SumMix([]any{"1", "5", "8", 8, 9, 9, 2, "3"}))
  fmt.Println(SumMix([]any{8, 0, 0, 8, 5, 7, 2, 3, 7, 8, 6, 7}))
}

func SumMix(arr []any) int {
  n:=0
  for _, v := range arr{
    switch v := v.(type){
    case int:
      n+=v
    case string:
      i, err := strconv.Atoi(v)
      if err!=nil {
        panic(err)
      }
      n+=i
    default:
      panic(fmt.Sprintf("unsupported type: %T", v))
    }
  }
  return n
}
