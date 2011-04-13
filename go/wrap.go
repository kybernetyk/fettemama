package main

//import "fmt"

func wordwrap(s string, maxlen int) string {
	ret := []byte(s)

	//fmt.Printf("s: %p\nr: %p\n", &s, ret)

	for i := 0; i < len(s); i++ {

		if i > 0 && i % maxlen == 0 {
			z := i
			if z < 0 {
				break
			}
			for {
				if s[z] == ' ' {
					ret[z] = '\n'
					break
				}
				z--
			}

		} else {
			ret[i] = s[i]
		}
	}

	return string(ret)
}
/*
func main() {
	s := `Ajajaj, jetzt macht das Bullshitbingo vor gar nichts mehr halt. Ich halte Test Driven Development ja (genau wie Pair Programming, etc) fuer ein Mittel um mittelmaessige Programmierer nicht in den Selbstmord zu treiben. "Hey, seht her. Ich kann mich keine 2 Minuten konzentrieren und der Debugger ist auch ein unbekanntes Wesen fuer mich! Aber ich kann mich Programmierer schimpfen, weil ich 75% meiner Zeit damit verbringe Testmodule zu schreiben".`
	x := wordwrap(s, 40)

	fmt.Println(x)

}*/
