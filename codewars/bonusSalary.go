package main

import (
	"fmt"
	"strconv"
)

/*
  if bonus is true, the salary should be multiplied by 10.
  If bonus is false, receive only his stated salary.
  Return with "£" = "\u00A3"
*/

func main() {
    fmt.Println(BonusTime(100, false))
    fmt.Println(BonusTime(9, false))
    fmt.Println(BonusTime(55000, false))
    fmt.Println(BonusTime(100, true))
    fmt.Println(BonusTime(14000, true))
}

func BonusTime(salary int, bonus bool) string {
	if bonus {
	  salary = salary * 10
	}
	return fmt.Sprintf("£%d", salary)
  }

func BonusTime2(salary int, bonus bool) string {
  var s int =0
  if bonus{
	s=salary*10
  } else {
	s=salary
  }

  st:=strconv.Itoa(s)

  return "\u00A3"+st
}