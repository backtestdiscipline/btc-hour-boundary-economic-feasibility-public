# Stage A economic feasibility report

Exploratory; previously exposed January 2025–September 2026 Binance Vision BTCUSDT spot archives. 918,720 rows, 15,311 eligible hours, 1 excluded hour, 92 UTC weeks. The first hour lacks its preceding candle. No fresh holdout is claimed.

## Movement in basis points

| Measure | Close-to-close minute 00 | Absolute open-to-close minute 00 |
|---|---:|---:|
| Mean | 4.544 | 4.545 |
| Median | 3.004 | 3.002 |
| Sample standard deviation | 5.580 | 5.582 |
| 90th percentile | 9.970 | 9.969 |
| 95th percentile | 13.797 | 13.807 |
| 99th percentile | 26.290 | 26.290 |
| Fraction above 5 bps | 30.93% | 30.92% |
| Fraction above 10 bps | 9.96% | 9.96% |
| Fraction above 15 bps | 4.21% | 4.20% |
| Fraction above 20 bps | 2.08% | 2.09% |
| Fraction above 30 bps | 0.66% | 0.66% |

Mean close-to-close minute 00: 4.544 [4.237, 4.880] bps; mean of minutes 01–59 within hour: 3.831 [3.560, 4.126]; difference: 0.714 [0.612, 0.821]. Brackets are two-sided 95% UTC-week cluster bootstrap intervals.

Signed long-only open-to-close mean: 0.040 [-0.123, 0.208] bps gross; hypothetical short-only: -0.040 [-0.208, 0.123] bps gross. Under primary 24 bps costs subtract 24 bps from each endpoint and mean. Neither is a fitted signal.

## Perfect-direction oracle sensitivity

Mean net bps [95% interval]; percentage of intervals with positive net. Fees per side and additional round-trip spread/slippage are separate assumptions.

| Fee per side | Extra 0 bps | Extra 2 bps | Extra 4 bps | Extra 8 bps |
|---:|---:|---:|---:|---:|
| 0 | 4.545 [4.237, 4.880]; 97.90% | 2.545 [2.237, 2.880]; 63.54% | 0.545 [0.237, 0.880]; 39.17% | -3.455 [-3.763, -3.120]; 15.41% |
| 2 | 0.545 [0.237, 0.880]; 39.17% | -1.455 [-1.763, -1.120]; 24.41% | -3.455 [-3.763, -3.120]; 15.41% | -7.455 [-7.763, -7.120]; 6.82% |
| 5 | -5.455 [-5.763, -5.120]; 9.96% | -7.455 [-7.763, -7.120]; 6.82% | -9.455 [-9.763, -9.120]; 4.83% | -13.455 [-13.763, -13.120]; 2.74% |
| 10 | -15.455 [-15.763, -15.120]; 2.09% | -17.455 [-17.763, -17.120]; 1.65% | -19.455 [-19.763, -19.120]; 1.31% | -23.455 [-23.763, -23.120]; 0.81% |

**NOT_ECONOMICALLY_PROMISING**. The primary 10 bps-per-side fee plus 4 bps extra gives -19.455 [-19.763, -19.120] bps per hour, wholly below zero. Do not pursue this particular one-minute directional approach under these assumptions. This does not rule out other volatility or execution-timing applications.

The oracle knows direction in advance and cannot be implemented. Open and close are recorded candle prices, not guaranteed fills; short exposure requires borrowing or a different instrument. Results are not a universal upper bound on intraminute approaches. Costs are hypothetical; 10 bps per side is consistent with the published regular-user [Binance spot schedule](https://www.binance.com/en/fee/trading), without assuming an individual tier or discount. Spread and slippage are additional and not measured. Bootstrap inference assumes week clusters sufficiently representative and independent; nonstationarity and the single venue limit generalization. Previously exposed periods cannot provide fresh validation.

Source: Binance Vision, BTCUSDT spot 1-minute monthly archives; [dataset terms](https://github.com/binance/binance-public-data/blob/master/TERMS_AND_CONDITIONS.md). Source archives are retained locally and excluded from Git. The frozen protocol precedes analysis. Platform integrity and semantic source replay are distinct checks; semantic replay uses the same calculation implementation, supplemented by independent synthetic examples. The earlier sealed run remains visible; a subsequent code change added the first-minute close-to-close confidence interval without changing the frozen method.
