# Texas Hold'em Poker — REST API

This repository contains a Go backend implementing a Texas Hold'em hand evaluator, comparator and Monte Carlo simulator, plus a minimal Flutter client and a tiny static web UI for quick testing.

## What this API does

- Evaluate a 7-card (2 hole + up to 5 community) poker hand and return its category and tiebreaker ranks.
- Compare two hands and return the best 5-card combinations and the winner.
- Run Monte Carlo simulations to estimate win/tie/loss percentages for a hero hand against random opponents.

## Endpoints

- `GET /health`
  - Returns `200 OK` and body `OK` when the service is healthy.

- `POST /evaluate`
  - Request JSON: `{ "hole": ["As","Kd"], "community": ["2d","7c","Jd"] }` (community optional)
  - Response JSON: `{ "rank": { "category": <int>, "ranks": [<int>, ...] } }`

- `POST /compare`
  - Request JSON: `{ "hand1": {"hole": [...], "community": [...]}, "hand2": {...} }`
  - Response JSON: `{ "winner": "hand1"|"hand2"|"tie", "hand1": {"best": ["As","Kd",...], "rank": {..}}, "hand2": {...} }`

- `POST /simulate`
  - Request JSON (hero vs random opponents):
    ```json
    {
      "hero": { "hole": ["As","Kd"] },
      "community_known": ["2d","7c","Jd"],
      "num_players": 2,
      "iterations": 1000,
      "concurrency": 4
    }
    ```
  - Response JSON:
    `{ "wins": <int>, "ties": <int>, "losses": <int>, "iterations": <int>, "winPct": <float>, "tiePct": <float>, "lossPct": <float> }`

## Hand ranking system

The evaluator returns a `category` integer and a `ranks` array used for tie-breaking. Categories follow this mapping (explicitly required):

0 = High Card
1 = One Pair
2 = Two Pair
3 = Three of a Kind
4 = Straight
5 = Flush
6 = Full House
7 = Four of a Kind
8 = Straight Flush

Card rank integers: `2..14` where `14` = Ace, `13` = King, ..., `2` = Two.

How to interpret `ranks`:

- `ranks` is an array of card rank integers ordered by tiebreaker significance. Compare two hands lexicographically by `category` first, then by `ranks` from left to right.

Example: `{ "category": 1, "ranks": [14, 13, 11] }` means a One Pair of Aces (pair of Aces), with kickers King (13), Jack (11).

## Example requests & responses

1) Evaluate example

Request:
```bash
curl -s -X POST https://texas-hold-em-poker.onrender.com/evaluate \
  -H "Content-Type: application/json" \
  -d '{"hole":["As","Kd"],"community":["2d","7c","Jd"]}'
```

Response:
```json
{
  "rank": { "category": 0, "ranks": [14, 13, 11, 7, 2] }
}
```

2) Compare example

Request:
```bash
curl -s -X POST https://texas-hold-em-poker.onrender.com/compare \
  -H "Content-Type: application/json" \
  -d '{"hand1":{"hole":["As","Ah"],"community":["Kd","Qd","Jd","2s","3c"]},"hand2":{"hole":["Ks","Kh"],"community":["Kd","Qd","Jd","2s","3c"]}}'
```

Response:
```json
{
  "winner": "hand1",
  "hand1": { "best": ["As","Ah","Kd","Qd","Jd"], "rank": { "category": 6, "ranks": [14, 13, 12] } },
  "hand2": { "best": ["Ks","Kh","Kd","Qd","Jd"], "rank": { "category": 6, "ranks": [13, 13, 12] } }
}
```

3) Simulate example

Request:
```bash
curl -s -X POST https://texas-hold-em-poker.onrender.com/simulate \
  -H "Content-Type: application/json" \
  -d '{"hero":{"hole":["As","Kd"]},"community_known":["2d","7c","Jd"],"num_players":2,"iterations":2000,"concurrency":4}'
```

Response (example):
```json
{
  "wins": 860,
  "ties": 40,
  "losses": 1100,
  "iterations": 2000,
  "winPct": 0.43,
  "tiePct": 0.02,
  "lossPct": 0.55
}
```

## Notes & troubleshooting

- The API accepts card strings in either `RankSuit` (e.g. `As`, `Td`, `10h`) or `SuitRank` (e.g. `S14`, `H10`) formats; ranks map `T=10, J=11, Q=12, K=13, A=14`.
- If you deploy frontend and API under the same domain the static web UI included here will call the API directly.
- For long-running simulations, reduce `iterations` or increase `concurrency` to tune performance.

## License
MIT
