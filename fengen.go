package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/amanjpro/chess"
	"github.com/amanjpro/zahak/engine"
	"github.com/amanjpro/zahak/evaluation"
	"github.com/amanjpro/zahak/search"
)

type (
	threadData struct {
		cache     *engine.Cache
		pawncache *evaluation.PawnCache
		runner    *search.Runner
		e         *search.Engine
	}
)

var ()

func init() {
}

func main() {
	lflag := flag.Int("limit", 0, "Maximum allowed difference between Quiescence Search result and Static Evaluation, the bigger it is the more tactical positions are included")
	iflag := flag.String("input", "", "Comma separated set of paths to PGN files")
	oflag := flag.String("output", "", "Directory to write produced FENs in")
	tflag := flag.Int("threas", runtime.NumCPU(), "Number of parallelism")
	flag.Parse()

	limit := int16(*lflag)
	inputPaths := strings.Split(*iflag, ",")
	if len(*iflag) == 0 || *iflag == "" {
		fmt.Println("At least the path of one PGN file is expected, none was given")
		os.Exit(1)
	}

	if len(*oflag) == 0 || *oflag == "" {
		fmt.Println("Output directory must be specified")
		os.Exit(1)
	}

	run(limit, inputPaths, *oflag, *tflag)
}

func run(limit int16, paths []string, outputDir string, threads int) {
	files := make([]*os.File, len(paths))
	channels := make([]chan *chess.Game, threads)
	outputs := make([]*bufio.Writer, threads)
	answer := make(chan int)
	for i := 0; i < threads; i++ {
		channels[i] = make(chan *chess.Game)
		f, err := os.Create(fmt.Sprintf("%s%cpart-%d.epd", outputDir, os.PathSeparator, i+1))
		if err != nil {
			panic(err)
		}
		outputs[i] = bufio.NewWriter(f)
		defer f.Sync()
		defer f.Close()

		c := engine.NewCache(1)
		pc := evaluation.NewPawnCache(2)
		r := search.NewRunner(c, pc, 1)
		t := threadData{
			cache:     c,
			pawncache: pc,
			runner:    r,
			e:         r.Engines[0],
		}
		t.runner.AddTimeManager(search.NewTimeManager(time.Now(), search.MAX_TIME, true, 0, 0, false))
		go t.process(limit, outputs[i], channels[i], answer)
	}

	nextThread := 0
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
			channels[nextThread] <- game
			nextThread = (nextThread + 1) % threads
		}
	}
	for i := 0; i < threads; i++ {
		close(channels[i])
	}
	fenCounter := 0
	for i := 0; i < threads; i++ {
		fenCounter += <-answer
	}
	fmt.Fprintf(os.Stderr, "Wrote %d FENs\n", fenCounter)
}

func (t *threadData) process(limit int16, output *bufio.Writer, games chan *chess.Game, answer chan int) {
	fenCounter := 0
	for game := range games {
		fenCounter += t.extractFens(game, limit, output)
	}
	answer <- fenCounter
}

func (t *threadData) extractFens(game *chess.Game, limit int16, out *bufio.Writer) int {
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
		t.e.Position = g.Position()
		if t.e.Position.IsInCheck() {
			continue // A position is in check? ignore it
		}
		t.runner.ClearForSearch()
		t.e.ClearForSearch()
		seval := evaluation.Evaluate(t.e.Position, t.pawncache)
		t.e.SetStaticEvals(0, seval)
		qeval := t.e.Quiescence(-engine.MAX_INT, engine.MAX_INT, 0)
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
			line := fmt.Sprintf("%s;score:%d;eval:%d;qs:%d;outcome:%s\n", fen, int(score*100), seval, qeval, outcome)
			_, err = out.WriteString(line)
			if err != nil {
				panic(err)
			}
			fenCounter += 1
		}
	}

	out.Flush()

	return fenCounter
}

func abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}
