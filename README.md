# FenGen

A simple utility to extract FENs from PGN databases. It uses
[Zahak](https://github.com/amanjpro/zahak) engine to decide wheather to keep a
FEN, or if it is too tactical discards it. Tactial positions are defined as:

- `ABS(Quiescence Eval - Static Eval)` is bigger than a threashold (configurable).
- Position is not in check
- Position is not a checkmate

This utility is useful for extracting FENs from PGN databases to be used for
tuning Hand-Crafted Evaluation terms, as well as training NNUE networks.

## Output Format

FenGen writes the generated FENs into stdout, and the format is as follows:

```
$ ./fengen -limit 40 -input ./self-play-1.pgn -output /tmp
$ cat /tmp/part-1.epd | head
1n1qkb1r/r2b1p2/4pn1p/p2p4/2pP1BpP/2P1P3/1PN2PP1/RN1QKB1R b KQk - 1 11;score:-31;eval:25;qs:25;outcome:0.5
r5kr/1p3ppp/1p6/1P1p4/3P2P1/qP5P/3NPP2/2RQK2R w K - 0 16;score:459;eval:449;qs:449;outcome:1.0
8/1r2p1k1/2Nnn1pp/R7/8/5NPP/5PK1/8 b - - 6 38;score:32;eval:12;qs:12;outcome:0.5
r3k1r1/1p6/2p1pb1p/2P2p2/p2P1P2/4P1P1/PP1B1K1P/1R2R3 w q - 2 27;score:188;eval:164;qs:164;outcome:1.0
6R1/pp6/1k4P1/r7/5K2/2p4P/8/8 w - - 2 54;score:-94;eval:-88;qs:-88;outcome:0.0
8/5k2/8/4p2p/4P3/2r2PP1/2p2K2/2R5 w - - 0 53;score:-123;eval:-112;qs:-112;outcome:0.5
8/3k1n2/1K5p/7R/6P1/8/8/8 b - - 3 55;score:315;eval:306;qs:306;outcome:0.5
8/6R1/8/5K2/5p2/2r2k2/8/8 w - - 4 67;score:-193;eval:-154;qs:-154;outcome:0.0
r2qkbr1/ppp1pb2/2n2n2/3pP2p/P2P2p1/2PQ2P1/1P1N1PB1/R1B1K1NR b KQq - 0 9;score:87;eval:-1;qs:-1;outcome:0.5
5rk1/R4pp1/5n1p/1r6/1bqp1P2/4PK1P/3P2P1/1N2Q1NR w - - 0 23;score:-98;eval:-145;qs:-145;outcome:0.0
```

- The first part is the FEN
- `score` is the score noted in the PGN comment of the move
- `eval` is the static eval of the position
- `qs` is the quiescence search result
- `outcome` is the game output, (0.0 white loses, 1.white wins, 0.5 is a draw).

All of `score`, `eval` and `qs` are from white's point of view.

## Usage

```
$ ./fengen -help
Usage of ./fengen:
  -input string
    Comma separated set of paths to PGN files
  -limit int
    Maximum allowed difference between Quiescence Search result and Static Evaluation, the bigger it is the more tactical positions are included
  -output string
    Directory to write produced FENs in
  -threas int
    Number of parallelism (default 8)
```
