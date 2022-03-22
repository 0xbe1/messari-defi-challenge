# Messari Defi Challenge

Challenge writeup: https://messari.notion.site/Messari-DeFi-Challenge-rev-03-17-2022-c5c6184e88dd44eab101be1f179a3ee0

GraphQL playground: https://thegraph.com/hosted-service/subgraph/uniswap/uniswap-v3

## Example

For liquidity pool ID: 0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8 (USDC/ETH), **US$356,129** in trading fees were collected on 16th March 2022.

Based on the total liquidity in the pool of US$391,636,206, this means that liquidity providers get to earn US$0.00090934 for every US$1 deposited into the pool.

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

- `id`: liquidity pool ID
- `feesUSD`: trading fee
- `tvlUSD`: total liquidity in the pool

To fetch data:

```gql
{
  poolDayDatas(where: {date_gte: 1637388800, date_lte: 1647930977}) {
    date
    pool {
      id
    }
    feesUSD
    tvlUSD
  }
}
```
