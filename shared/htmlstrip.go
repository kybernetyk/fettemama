package main

/*import (
	"fmt"
)
*/


//I hope I got the slices right and am not copying strings aroung :]
func htmlstrip(s string) string {
	l := len(s)
	sl := []byte(s)

	ts := make([]byte, 0, l) //our return slice

	tmp := 0
	for i, c := range s {
		if c == '<' {
			idx := i
			if idx < 0 {
				idx = 0
			}
			ts = append(ts, sl[tmp:idx]...)
			tmp = idx
		}
		if c == '>' {
			tmp = i + 1
			if tmp > l {
				tmp = l
			}
		}
	}
	ts = append(ts, sl[tmp:len(s)]...)

	return string(ts)
}
/*
var x = `Living the future:
<blockquote>Die US-Marine hat erstmals einen Hochenergie-Laser auf See abgefeuert - und bei dem Experiment ein kleines Boot in Brand gesetzt. Schiffe sollen sich künftig mit solchen Energiewaffen verteidigen</blockquote>
(<a href="http://www.spiegel.de/wissenschaft/technik/0,1518,756514,00.html">via</a>)
`

func main() {
	strs := []string{"<b>ich bin behinat</b>.",
		"<a href='http://www.de'>lol</a>",
		"das ist ein <p>test</p>",
		x,
		"omg <b>d</b> lol<br>o",
		"I <3 u all!",
		"I <3 u <b>all!</b>",
		"I <b>love</b> thats <3 u all!!!!1"}
	for _, s := range strs {
		fmt.Println([]byte(htmlstrip(s)))
	}
}*/
