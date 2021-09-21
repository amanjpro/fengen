package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amanjpro/chess"
	"github.com/amanjpro/zahak/engine"
	"github.com/amanjpro/zahak/evaluation"
	"github.com/amanjpro/zahak/search"
)

var (
	cache     = engine.NewCache(1)
	pawncache = evaluation.NewPawnCache(2)
	runner    = search.NewRunner(cache, pawncache, 1)
	e         = runner.Engines[0]
)

func init() {
	runner.AddTimeManager(search.NewTimeManager(time.Now(), search.MAX_TIME, true, 0, 0, false))
}

func main() {
	lflag := flag.Int("limit", 0, "Maximum allowed difference between Quiescence Search result and Static Evaluation, the bigger it is the more tactical positions are included")
	pflag := flag.String("paths", "", "Comma separated set of paths to PGN files")
	flag.Parse()

	limit := int16(*lflag)
	paths := strings.Split(*pflag, ",")
	if len(paths) == 0 || *pflag == "" {
		panic("At least the path of one PGN file is expected, none was given")
	}

	process(limit, paths)
}

func process(limit int16, paths []string) {
	files := make([]*os.File, len(paths))
	fenCounter := 0

	for i, p := range paths {
		f, err := os.Open(p)
		files[i] = f
		if err != nil {
			panic(err)
		}
		defer files[i].Close()

		scanner := chess.NewScanner(f)
		for scanner.Scan() {
			game := scanner.Next()
			fenCounter += extractFens(game, limit)
		}
	}
	fmt.Fprintf(os.Stderr, "Wrote %d FENs\n", fenCounter)
}

func extractFens(game *chess.Game, limit int16) int {
	comments := game.Comments()
	var outcome string
	if game.Outcome() == chess.WhiteWon {
		outcome = "1.0"
	} else if game.Outcome() == chess.BlackWon {
		outcome = "0.0"
	} else if game.Outcome() == chess.Draw {
		outcome = "0.5"
	} else {
		return 0 // no outcome? go to the next game
	}
	fenCounter := 0
	for i, pos := range game.Positions() {
		if i == 0 {
			continue // Not intersted in the startpos
		}
		if i == len(game.Positions()) && game.Method() == chess.Checkmate {
			continue // ignore checkamte positions
		}
		fen := pos.String()
		g := engine.FromFen(fen)
		e.Position = g.Position()
		if e.Position.IsInCheck() {
			continue // A position is in check? ignore it
		}
		runner.ClearForSearch()
		e.ClearForSearch()
		seval := evaluation.Evaluate(e.Position, pawncache)
		e.SetStaticEvals(0, seval)
		qeval := e.Quiescence(-engine.MAX_INT, engine.MAX_INT, 0)
		if abs16(seval-qeval) <= limit {
			tokens := strings.Split(comments[i-1][0], " ")
			scoreStr := strings.Split(tokens[0], "/")[0]
			score, err := strconv.ParseFloat(scoreStr, 64)
			if err != nil {
				if scoreStr == "book" || strings.Contains(scoreStr, "M") {
					continue // Not interested in near mate positions, or book moves
				}
				panic(err)
			}
			if math.Abs(score) > 2000 {
				continue // Not interested decided positions
			}
			// Last move was black's move, i.e. black has reported the scores
			// To convert it into white's point of view, we need to negate the scores
			if pos.Turn() == chess.White {
				score *= -1
			} else {
				seval *= -1
				qeval *= -1
			}
			fmt.Printf("%s;score:%d;eval:%d;qs:%d;outcome:%s\n", fen, int(score*100), seval, qeval, outcome)
			fenCounter += 1
		}
	}

	return fenCounter
}

func abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
