package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

// OMap is the struct that can be serialized as JSON
// containing all the information to create a new ConsMap
type OMap struct {
	ResetT    time.Duration     `json:"resetTime"`
	LastReset time.Time         `json:"lastReset"`
	UserCons  map[string]uint64 `json:"userCons"`
}

func main() {
	var fp string
	flag.StringVar(&fp, "f", "",
		"Path to JSON serialized map with keys of type string")
	flag.Parse()
	rg := regexp.MustCompile("[[:upper:]]")
	f, e := os.Open(fp)
	mp := new(OMap)
	if e == nil {
		d := json.NewDecoder(f)
		e = d.Decode(mp)
		f.Close()
	}
	var w *os.File
	var rmp *OMap
	if e == nil {
		rmp = &OMap{
			LastReset: mp.LastReset,
			ResetT:    mp.ResetT,
			UserCons:  make(map[string]uint64),
		}
		for k, v := range mp.UserCons {
			if !rg.MatchString(k) {
				rmp.UserCons[k] = v
			}
		}
		w, e = os.Create(fp)
	}
	if e == nil {
		n := json.NewEncoder(w)
		n.SetIndent("	", "")
		e = n.Encode(rmp)
		w.Close()
	}
	ex := 0
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		ex = 1
	}
	os.Exit(ex)
}
