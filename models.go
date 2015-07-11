package main

import (
	"time"
)

type Book struct {
	Id      string `gorethink:"id,omitempty"`
	Title   string
	Author  string
	Iscale  float64
	Lscale  float64
	Genre   []string
	Created time.Time
}
