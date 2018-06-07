package app

// Season is a container for season related information
type Season struct {
	Name     string
	Number   int
	Episodes []*Episode
}

func (s Season) Len() int {
	return len(s.Episodes)
}
func (s Season) Swap(i, j int) {
	s.Episodes[i], s.Episodes[j] = s.Episodes[j], s.Episodes[i]
}
func (s Season) Less(i, j int) bool {
	return s.Episodes[i].Number < s.Episodes[j].Number
}
