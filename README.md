# crypto-curry-print

List crypto currencies by last coins listed @ coinmarketcap

## Dependencies

`go get github.com/tyler-sommer/stick`

`go get modernc.org/sqlite`

## Building

`go build`

## Configuration

place your CoinMarketCap apiKey in `.env.json`
place your BitQuery apiKey in `.env.json`

## Usage

`./cryptocurryprint 8880`

Explore last added cryptocurrencies @ `http://127.0.0.1:8880`
