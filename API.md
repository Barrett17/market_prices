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

Example:

curl -X GET http://127.0.0.1:6060/api/ticker/GBP

{
"ticker" : "GBP",
"rate" : 1.1216421,
"weeklyrate" : 1.1246789,
"buy" : false
}
