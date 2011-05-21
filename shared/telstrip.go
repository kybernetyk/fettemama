package main

//strip evil console command codes out ...
func telstrip(s string) string {
	ts := make([]int, 0, len(s)) //our return slice

	for _, c := range s {
		if c < ' ' {
			continue
		}
		ts = append(ts, c)
	}

	return string(ts)
}
