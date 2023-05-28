package main

type AlphabetCounter struct {
	letterCount         [26]int
	maxRecurLetter      string
	maxRecurLetterCount int
}

func (l *AlphabetCounter) GetHighestCountLetter() (string, int) {
	return (*l).maxRecurLetter, (l).maxRecurLetterCount
}

func (l *AlphabetCounter) AddLetter(val string) {
	asc11 := int(rune(val[0])) - 65
	(*l).letterCount[asc11]++

	if (*l).maxRecurLetter == "" || (*l).maxRecurLetter == val {
		(*l).maxRecurLetter = val
		(*l).maxRecurLetterCount = (*l).letterCount[asc11]
	} else {
		currentLetterCount := (*l).letterCount[asc11]
		if currentLetterCount > (*l).maxRecurLetterCount {
			(*l).maxRecurLetter = val
			(*l).maxRecurLetterCount = (*l).letterCount[asc11]
		}
	}
}