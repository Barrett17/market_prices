/api/ticker/{pair}
[GET]

Get price in euro for the pair selected. Currently available pairs:

USD-EUR
GBP-EUR

Returns (JSON):

- ticker

The ticker used to make the request.

- rate

The current rate of the trading pair.

- weekrate

The last week average rate.

- prediction

The prediction for the next week (buy/sell)

