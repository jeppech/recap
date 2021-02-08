# recap-go
This library will map NAMED regex capture groups, to a tagged struct

```golang
package main

import "github.com/jeppech/recap-go"

type Person struct {
  Name string   `recap:"name"`
  Age int       `recap:"age"`
  Public bool   `recap:"public;default=false"`
  Apple bool    `recap:"fruit;contains=apple"`
  Orange bool   `recap:"fruit;contains=orange"`
  Banana bool   `recap:"fruit;contains=banana"`
}

func main() {
  pers := Person{}
  recap.Parse(&pers, `(?P<name>[A-Za-z]+) (?P<age>\d+) (?:(?P<public>(true|false)))?\W?\[(?P<fruit>.*)\]`, "Jeppe 31 [orange banana]")
  fmt.Printf("%+v\n", pers)
}

// {Name:Jeppe Age:31 Public:false Apple:false Orange:true Banana:true}
```