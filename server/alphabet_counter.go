package main

type AlphabetCounter struct {
	letterCount          [26]int64
	maxRecurLetter       string
	maxRecurLetterCount  int64
	totalLettersReceived int64
}

func (l *AlphabetCounter) GetCounterStats() (string, int64, int64) {
	return (*l).maxRecurLetter, (*l).maxRecurLetterCount, (*l).totalLettersReceived
}

func (l *AlphabetCounter) AddLetter(val string) {
	asc11 := int(rune(val[0])) - 65
	(*l).letterCount[asc11]++
	(*l).totalLettersReceived++

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