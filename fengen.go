package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/amanjpro/zahak/engine"
	"github.com/amanjpro/zahak/evaluation"
	"github.com/amanjpro/zahak/search"
	"github.com/notnil/chess"
)

func main() {
	lflag := flag.Int("limit", 0, "Maximum allowed difference between Quiescence Search result and Static Evaluation, the bigger it is the more tactical positions are included")
	pflag := flag.String("paths", "", "Comma separated set of paths to PGN files")
	flag.Parse()

	limit := int16(*lflag)
	paths := strings.Split(*pflag, "\n")
	if len(paths) == 0 || *pflag == "" {
		panic("At least the path of one PGN file is expected, none was given")
	}
	files := make([]*os.File, len(paths))

	for i, p := range paths {
		f, err := os.Open(p)
		files[i] = f
		if err != nil {
			panic(err)
		}
		defer files[i].Close()

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
				if abs16(seval-qeval) <= limit {
					fmt.Printf("%s;static:%d;qs:%d\n", fen, seval, qeval)
				}
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
