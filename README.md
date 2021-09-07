# FenGen

A simple utility to extract FENs from PGN databases. It uses Zahak engine to
decide wheather to keep a FEN or if it is too tactical discards it. Tactial
position is defined as: `ABS(Quiescence Eval - Static Eval)` is bigger than the
a threashold.

This utility is useful for extracting FENs from PGN to be used for tuning
Hand-Crafted Evaluation terms, as well as training NNUE networks.

The output is as follows:

```
$ ./fengen -limit 40 -paths ./self-play-1.pgn | head
rnbqkb1r/2pp1p1p/1p2pn2/p5p1/PP6/N6N/1BPPPPPP/R2QKB1R w KQkq - 2 2;score:-0.450000;eval:-6;qs:5,outcome:0.0
rnbqkb1r/2pp1p1p/1p2pn2/P5p1/P7/N6N/1BPPPPPP/R2QKB1R b KQkq - 0 2;score:-0.550000;eval:-19;qs:-57,outcome:0.0
1nbqkb1r/2pp1p1p/1p2pn2/r5p1/P3P3/N6N/1BPP1PPP/R2QKB1R b KQk e3 0 3;score:-0.350000;eval:-31;qs:0,outcome:0.0
1nbqkb1r/2pp1p1p/1p2pn2/r7/P3P1p1/N6N/1BPP1PPP/R2QKB1R w KQk - 0 4;score:-0.380000;eval:-35;qs:0,outcome:0.0
1nbqkb1r/2pp1p1p/1p2pn2/r7/P3P1p1/N7/1BPP1PPP/R2QKBNR b KQk - 1 4;score:-0.090000;eval:-65;qs:-70,outcome:0.0
1nbqkb1r/2p2p1p/1p2p3/r2pP3/P3n1Q1/N7/1BPP1PPP/R3KBNR b KQk - 0 6;score:0.200000;eval:-9;qs:-29,outcome:0.0
2bqkb1r/2pn1p1p/1p2p3/r2pP3/P3n1Q1/N7/1BPP1PPP/R3KBNR w KQk - 1 7;score:0.020000;eval:-38;qs:0,outcome:0.0
2bqkb1r/2pn1p1p/1p2p3/r2pP3/P1P1n1Q1/N7/1B1P1PPP/R3KBNR b KQk c3 0 7;score:0.220000;eval:-38;qs:0,outcome:0.0
2bqk2r/2pn1p1p/1p2p3/r2pP3/PbP1n1Q1/N4N2/1B1P1PPP/R3KB1R b KQk - 2 8;score:-0.970000;eval:12;qs:0,outcome:0.0
2bqkn1r/2p2p1p/1p2p3/r2pP3/PbP1n1Q1/N4N2/1B1P1PPP/R3KB1R w KQk - 3 9;score:-0.940000;eval:27;qs:61,outcome:0.0
```

The first part is the FEN, score is the score noted in the PGN comment of the
move, the eval is the static eval, qs is the quiescence search result, outcome
is the winner (0.0 black wins, 1.0 white wins, 0.5 is a draw). All of score,
eval and quiescence score are from white's point of view.

Usage:

```
$ ./fengen -help
Usage of ./fengen:
  -limit int
    Maximum allowed difference between Quiescence Search result and Static Evaluation, the bigger it is the more tactical positions are included
  -paths string
    Comma separated set of paths to PGN files
```
