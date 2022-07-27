package main

import "time"

type Doctor struct {
	ID      int64
	Name    string
	Age     int
	Sex     int // 1 for male , 2 for female
	AddTime time.Time
}
