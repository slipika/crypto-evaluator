
# crypto-evaluator

crypto-evaluator is a command line based utility that takes in a USD amount as holdings, and
calculates the 70/30 split for 2 given crypto currencies. This tool uses COINBASE api to fetch exchange rates for different crypto currency. Each rate is how much of that crypto currency you would get
for 1 dollar. These rates are used to calculate the 70/30 split of USD amount for the two given currencies.

## Getting started

### Build & Install

#### 1.1 Clone this repo
#### 1.2 Build the CLI

``` 
make build 
```
### Run the tool
```
 ./crypto-evaluator usd-amount currency-70-split currency-30-split 
 ```
 ### Usage
#### 1. Example Input
```
./crypto-evaluator 200 ZXD ZRX
```

#### 2. Example output
```
140=>52833.4923
60=>168.4177
```
The first line represents the seventy split that is for a total of 200 dollars holding is 140 dollars is the seventy split  with 52833.4923 in ZXD crypto currency. The second line represents the 30% holdings in ZRX currency.

### Run tests
```
make test
make test_coverage
```
