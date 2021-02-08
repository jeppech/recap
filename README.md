# Recap
This library will map NAMED regex capture groups, to a tagged struct

```golang
package main

import "github.com/jeppech/recap-go"

type Person struct {
	Name string		`recap:"name"`
	Age int				`recap:"age"`
	Public bool		`recap:"public;default=false"`
}

func main() {
	pers := Person{}
	recap.Parse(pers, `(?P<name>) (?P<age>) (?P<public>)`, "Jeppe 31")
	fmt.Printf("%+v\n", pers)
}
```