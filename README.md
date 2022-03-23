# Messari Defi Challenge

Challenge writeup: https://messari.notion.site/Messari-DeFi-Challenge-rev-03-17-2022-c5c6184e88dd44eab101be1f179a3ee0

GraphQL playground: https://thegraph.com/hosted-service/subgraph/uniswap/uniswap-v3

## So, how to get the data?

According to the writeup:

> For liquidity pool ID: 0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8 (USDC/ETH), **US$356,129** in trading fees were collected on 16th March 2022.
>
> Based on the total liquidity in the pool of US$391,636,206, this means that liquidity providers get to earn US$0.00090934 for every US$1 deposited into the pool.

To fetch example data:

```gql
{
  poolDayDatas(where: {date: 1647388800, pool: "0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8"}) {
    date
    pool {
      id
      token0 {
        name
      }
      token1 {
        name
      }
    }
    feesUSD
    tvlUSD
  }
}
```

Result:

```json
{
  "data": {
    "poolDayDatas": [
      {
        "date": 1647388800,
        "pool": {
          "id": "0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8",
          "token0": {
            "name": "USD Coin"
          },
          "token1": {
            "name": "Wrapped Ether"
          }
        },
        "feesUSD": "356129.9553516688861767489783604968",
        "tvlUSD": "391636206.112549300335557572015702"
      }
    ]
  }
}
```

We learn that:

- `pool.id`: liquidity pool ID
- `feesUSD`: trading fee
- `tvlUSD`: total liquidity in the pool

## Next, how to tackle the challenge?

Steps:

1. Fetch all records within the given timeframe
1. For each pool, **earning rate of each day** = feesUSD / tvlUSD, **earning rate** = sum(earning rate of each day) over the timeframe
1. Find the pool with largest earning rate

See main.go. To run the program, simply `make run`.

## Result analysis

The program prints:

```
0x7845cfd7acb64e988988f0eeec47ec84c4fb0021
9.32681182538145e+15
```

The 2nd line looks odd - why is it so large? Investigate the pool ID, we find a record with fees=28.51685188526597089835106461853531 and tvl=0.000000000000003057513373183090360444999031674324. This contributes to the giant earning rate.
