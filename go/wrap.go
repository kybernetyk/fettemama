package main

import "strings"

//why can't you have char*, Go?
func wordwrap(s string, maxlen int) string {
	ret := make([]byte, 0, len(s))

	indent := false
	p := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			p = 0
		}
		if s[i] == 000 {
			indent = true
			ret = append(ret, '\003')
			continue
		}
		if s[i] == 001 {
			indent = false
			continue
		}

		if p > 0 && p%maxlen == 0 {
			z := i
			if z < 0 {
				break
			}
			for {
				if s[z] == ' ' {
					p = 0
					d := i - z
					i = z

					ret = ret[0 : len(ret)-d]
					ret = append(ret, '\n')
					if indent {
						ret = append(ret, '\003')
					}
					break
				}
				z--
			}

		} else {
			ret = append(ret, s[i])
			if s[i] == '\n' && indent {
				ret = append(ret, '\003')
			}
		}
		p++
	}

	str := strings.Replace(string(ret), "\003", "  > ", -1)

	return str
}
/*
func main() {
	s := `Ajajaj, jetzt macht das Bullshitbingo vor gar nichts mehr halt. Ich halte Test Driven Development ja (genau wie Pair Programming, etc) fuer ein Mittel um mittelmaessige Programmierer nicht in den Selbstmord zu treiben. "Hey, seht her. Ich kann mich keine 2 Minuten konzentrieren und der Debugger ist auch ein unbekanntes Wesen fuer mich! Aber ich kann mich Programmierer schimpfen, weil ich 75% meiner Zeit damit verbringe Testmodule zu schreiben".`
	x := wordwrap(s, 40)

	fmt.Println(x)

}*/
