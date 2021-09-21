# recap
This library will map NAMED regex capture groups, to a tagged struct

```golang
package main

import "github.com/jeppech/recap"

type Person struct {
  Name string   `recap:"name"`
  Age int       `recap:"age"`
  Public bool   `recap:"public;default=false"`
  Fruits Fruits
}

type Fruits struct {
  Apple bool    `recap:"fruit;contains=apple"`
  Orange bool   `recap:"fruit;contains=orange"`
  Banana bool   `recap:"fruit;contains=banana"`
}

func main() {
  pers := Person{}
	rx := regexp.MustCompile(`(?P<name>[A-Za-z]+) (?P<age>\d+) (?:(?P<public>(true|false)))?\W?\[(?P<fruit>.*)\]`)
  err, match := recap.Parse(&pers, rx, "Jeppe 31 [orange banana]")
	if err != nil {
		log.Fatal(err)
	}

	if (match) {
		fmt.Printf("%+v\n", pers)
	} else {
		// did not match
	}
}

// {Name:Jeppe Age:31 Public:false Fruits:{Apple:false Orange:true Banana:true}}
```