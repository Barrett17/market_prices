/api/ticker/{ticker}
[GET]

Get price in euro for the pair selected. Currently available pairs:

USD
GBP

Returns (JSON):

- ticker (string)

The ticker used to make the request.

- rate (float)

The current rate of the trading pair.

- weekrate (float)

The last week average rate.

- buy (bool)

The prediction for the next week (buy=true/sell=false)

