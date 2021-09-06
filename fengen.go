package main

import (
	"fmt"
	"os"
	"time"

	"github.com/amanjpro/zahak/engine"
	"github.com/amanjpro/zahak/evaluation"
	"github.com/amanjpro/zahak/search"
	"github.com/notnil/chess"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		panic(fmt.Sprintf("Exactly one argument is expected, the name of the source PGN file, %s", args))
	}
	f, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cache := engine.NewCache(1)
	pawncache := evaluation.NewPawnCache(2)
	runner := search.NewRunner(cache, pawncache, 1)
	runner.AddTimeManager(search.NewTimeManager(time.Now(), 838838292838383, true, 0, 0, false))
	e := runner.Engines[0]
	scanner := chess.NewScanner(f)
	for scanner.Scan() {
		game := scanner.Next()
		for _, pos := range game.Positions() {
			fen := pos.String()
			g := engine.FromFen(fen)
			e.Position = g.Position()
			runner.ClearForSearch()
			e.ClearForSearch()
			seval := evaluation.Evaluate(e.Position, pawncache)
			qeval := e.Quiescence(-engine.MAX_INT, engine.MAX_INT, 0)
			if abs16(seval-qeval) <= 50 {
				fmt.Println(fen, "|", qeval, "|", seval)
			}
		}
	}
}

func abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
