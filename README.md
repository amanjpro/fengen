# FenGen

A simple utility to extract FENs from PGN databases. It uses Zahak engine to
decide wheather to keep a FEN or if it is too tactical discards it.  Tactial
position is defined as: `ABS(Quiescence Eval - Static Eval)` is bigger than the
value of half a pawn
