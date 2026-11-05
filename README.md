# BTC hour-boundary economic feasibility

Can the first minute of a UTC hour support a directional BTC trade after costs? This exploratory screen used previously exposed Binance Vision BTCUSDT spot candles from January 2025 through September 2026. A hypothetical oracle knew each first minute's direction in advance and was credited with the candle open and close as fills. Under the frozen assumption of 10 basis points per side plus 4 bps additional round-trip costs, its mean net return was −19.455 bps per interval (two-sided 95% UTC-week bootstrap interval −19.763 to −19.120). The accepted classification was `NOT_ECONOMICALLY_PROMISING` for this specified approach and cost assumption. Neither the oracle nor its fills are executable guarantees; short exposure requires borrowing or another instrument. This is not a fresh holdout, a predictive signal, or a universal verdict on intraminute trading.

Read the [publication](https://backtestdiscipline.com/research/btc-hour-boundary-economic-feasibility/), [protocol](PROTOCOL.md) and [report](REPORT.md).

## Evidence and verification

The [frozen records](records/btc-hour-boundary-economic-feasibility/) specify the source identities, method and decision rule. Two complete sealed runs are retained: [first](.research-platform/runs/20261105T163942Z-1d113c529f09/) and [second](.research-platform/runs/20261105T164021Z-30a4b1ce8180/). The second adds a previously omitted first-minute close-to-close confidence interval without changing the frozen method. The [handoff](handoffs/20261105T164021Z-30a4b1ce8180.json) records both runs. The second run's result SHA-256 is `d7d67dcc6272783ed74aef8df05a678d65d8c224974eaf4d9554e590dfceea43`; semantic verification SHA-256 is `a023b0c8e6e907c42f1794331a62dad7df69082ba60569db748d1692980dbdcc`.

The recorded platform integrity check and study-owned semantic replay both passed. Integrity checks byte identities and bindings; semantic replay checks the scientific calculation from the bound sources using the same implementation, supplemented by independent synthetic examples. Neither is a third-party audit. Hashes identify bytes and do not by themselves establish chronology or validity. These are recorded results, not new verification by this snapshot.

The frozen protocol says malformed timestamps fail. The retained decoder instead excludes up to 1,000 invalid close timestamps; the sealed result records zero invalid close-timestamp rows, so this discrepancy did not affect the reported sample. The frozen protocol and implementation are retained unchanged.

## Code and reproducibility

The original study's retained Git HEAD was `bc1c1bec7499162674f1e198ca3be0769d4f049b`. This is a fresh public snapshot, not that repository's complete Git history. Frozen records, run plans, seals, declarations, results and handoff are preserved byte-for-byte; recorded commits and implementation hashes identify the original private history, which this snapshot does not supply. The [study code](study/) and [CLI](cmd/) are included for inspection; the second run's implementation added the close-to-close interval after the first run. Replaying the first run requires its recorded earlier implementation, which is not supplied by this snapshot.

The raw archives and separate private `research-platform` dependency are not included. The Go code imports that platform and `go.mod` refers to a sibling checkout. Consequently this snapshot alone cannot run its tests or perform platform verification. Full replay requires the matching platform implementation and exact source archives matching the recorded hashes; a publisher may correct or remove files. When an authorized compatible sibling platform is available, Go 1.27.2 and `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test ./...` run the study tests. This is partial inspection, not a self-contained rerun package.

## Attribution and licensing

Underlying candles: [Binance Vision BTCUSDT spot one-minute archives](https://data.binance.vision/?prefix=data/spot/monthly/klines/BTCUSDT/1m/). Backtest Discipline calculated signed and absolute log returns, cost-adjusted oracle returns and aggregates. The derived analytical results are offered under [CC BY-NC-SA 4.0](https://creativecommons.org/licenses/by-nc-sa/4.0/) with attribution to Binance Vision and Backtest Discipline, subject to the [Binance Vision Dataset Terms](https://github.com/binance/binance-public-data/blob/master/TERMS_AND_CONDITIONS.md) and general terms. This does not license raw candles, unrelated prose, or the private platform. Binance does not endorse the analysis. See [NOTICE](NOTICE).

The original Go code included in `study/` and `cmd/` is Copyright 2026 Mikhail Klemin, licensed under [Apache 2.0](LICENSE). This grant excludes the platform, data, and third-party material.
