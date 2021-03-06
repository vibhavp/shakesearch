package parse

import (
	"index/suffixarray"
	"log"
	"sort"
)

type Work struct {
	Title string
	Index *suffixarray.Index
	Start int
	End   int
}

var titles []string = []string{
	"THE SONNETS",
	"ALL’S WELL THAT ENDS WELL",
	"THE TRAGEDY OF ANTONY AND CLEOPATRA",
	"AS YOU LIKE IT",
	"THE COMEDY OF ERRORS",
	"THE TRAGEDY OF CORIOLANUS",
	"CYMBELINE",
	"THE TRAGEDY OF HAMLET, PRINCE OF DENMARK",
	"THE FIRST PART OF KING HENRY THE FOURTH",
	"THE SECOND PART OF KING HENRY THE FOURTH",
	"THE LIFE OF KING HENRY V",
	"THE FIRST PART OF HENRY THE SIXTH",
	"THE SECOND PART OF KING HENRY THE SIXTH",
	"THE THIRD PART OF KING HENRY THE SIXTH",
	"KING HENRY THE EIGHTH",
	"KING JOHN",
	"THE TRAGEDY OF JULIUS CAESAR",
	"THE TRAGEDY OF KING LEAR",
	"LOVE’S LABOUR’S LOST",
	"MACBETH",
	"MEASURE FOR MEASURE",
	"THE MERCHANT OF VENICE",
	"THE MERRY WIVES OF WINDSOR",
	"A MIDSUMMER NIGHT’S DREAM",
	"MUCH ADO ABOUT NOTHING",
	"OTHELLO, THE MOOR OF VENICE",
	"PERICLES, PRINCE OF TYRE",
	"KING RICHARD THE SECOND",
	"KING RICHARD THE THIRD",
	"THE TRAGEDY OF ROMEO AND JULIET",
	"THE TAMING OF THE SHREW",
	"THE TEMPEST",
	"THE LIFE OF TIMON OF ATHENS",
	"THE TRAGEDY OF TITUS ANDRONICUS",
	"THE HISTORY OF TROILUS AND CRESSIDA",
	"TWELFTH NIGHT: OR, WHAT YOU WILL",
	"THE TWO GENTLEMEN OF VERONA",
	"THE TWO NOBLE KINSMEN",
	"THE WINTER’S TALE",
	"A LOVER’S COMPLAINT",
	"THE PASSIONATE PILGRIM",
	"THE PHOENIX AND THE TURTLE",
	"THE RAPE OF LUCRECE",
	"VENUS AND ADONIS",
}

func GetWorks(index *suffixarray.Index) ([]Work, error) {
	allWorks := make([]Work, len(titles))

	for i, title := range titles {
		res := index.Lookup([]byte(title), -1)
		sort.Ints(res)
		log.Println("Finding", title)

		allWorks[i] = Work{
			Title: title,
			Start: res[0],
		}

		if i > 0 {
			allWorks[i-1].End = allWorks[i].Start - 1
			allWorks[i-1].Index = suffixarray.New(index.Bytes()[allWorks[i-1].Start:allWorks[i].Start])
		}
	}

	allWorks[len(titles)-1].End = len(index.Bytes()) - 1
	allWorks[len(titles)-1].Index = suffixarray.New(index.Bytes()[allWorks[len(titles)-1].Start:])

	return allWorks, nil

}
