# Stocker
Go command line tool for rebalancing stock portfolio. Currently only knows how
to pull current stock values from Vanguard.

## To use

1. Create a `holdings.json` file similar to the one below:
   ```javascript
   {
     "Holdings": {
       "VXUS": 177.487,
       "VTI": 145.163,
     },
     "TargetRatio": {
       "VXUS": 25,
       "VTI": 75
     }
   }
   ```
   This file says 75% of the total value should be invested in VTI, and 25% should
be invested in VXUS.
1. Run the tool with `go run *.go --holdings=holdings.json --amount=1234.56`,
   where `1234.56` is how much money you're adding to the fund.
