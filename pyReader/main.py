import time
import sys
import os
from threading import Thread
from multiprocessing import Process, Lock

import datetime as DT

from web3 import Web3
from web3.middleware import geth_poa_middleware  # Needed for Binance

from json import loads
from decimal import Decimal

ETHER = 10 ** 18
providers = [
    "https://bsc-dataseed.binance.org/",
    "https://bsc-dataseed1.defibit.io/",
    "https://bsc-dataseed1.ninicoin.io",
    "https://bsc-dataseed2.defibit.io/",
    "https://bsc-dataseed3.defibit.io/",
    "https://bsc-dataseed4.defibit.io/",
    "https://bsc-dataseed2.ninicoin.io/",
    "https://bsc-dataseed3.ninicoin.io/",
    "https://bsc-dataseed4.ninicoin.io/",
    "https://bsc-dataseed1.binance.org/",
    "https://bsc-dataseed2.binance.org/",
    "https://bsc-dataseed3.binance.org/",
    "https://bsc-dataseed4.binance.org/"
]
providersSuccess = [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ]

WBNB = '0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c'

web3 = Web3(Web3.HTTPProvider('https://bsc-dataseed1.binance.org:443'))
web3.middleware_onion.inject(geth_poa_middleware, layer=0)  # Again, this is needed for Binance, not Ethirium

ERC20_ABI = '[{"constant": true,"inputs": [],"name": "name","outputs": [{"name": "","type": "string"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "_spender","type": "address"},{"name": "_value","type": "uint256"}],"name": "approve","outputs": [{"name": "","type": "bool"}],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": true,"inputs": [],"name": "totalSupply","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "_from","type": "address"},{"name": "_to","type": "address"},{"name": "_value","type": "uint256"}],"name": "transferFrom","outputs": [{"name": "","type": "bool"}],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": true,"inputs": [],"name": "decimals","outputs": [{"name": "","type": "uint8"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": true,"inputs": [{"name": "_owner","type": "address"}],"name": "balanceOf","outputs": [{"name": "balance","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": true,"inputs": [],"name": "symbol","outputs": [{"name": "","type": "string"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "_to","type": "address"},{"name": "_value","type": "uint256"}],"name": "transfer","outputs": [{"name": "","type": "bool"}],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": true,"inputs": [{"name": "_owner","type": "address"},{"name": "_spender","type": "address"}],"name": "allowance","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"payable": true,"stateMutability": "payable","type": "fallback"},{"anonymous": false,"inputs": [{"indexed": true,"name": "owner","type": "address"},{"indexed": true,"name": "spender","type": "address"},{"indexed": false,"name": "value","type": "uint256"}],"name": "Approval","type": "event"},{"anonymous": false,"inputs": [{"indexed": true,"name": "from","type": "address"},{"indexed": true,"name": "to","type": "address"},{"indexed": false,"name": "value","type": "uint256"}],"name": "Transfer","type": "event"}]'

ABI = loads(
    '[{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"token0","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[],"name":"factory","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"address","name":"","type":"address"}],"name":"getPair","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"}]')

providerIndex = 0

priceData = []
dict = {}

def get_iso8601_from_timestamp(block_timestamp):
    return str(DT.datetime.utcfromtimestamp(block_timestamp)).replace(" ", "T") + ".000Z"


def get_price(token, decimals, pair_contract, is_reversed, is_price_in_peg):
    peg_reserve = 0
    token_reserve = 0
    (reserve0, reserve1, blockTimestampLast) = pair_contract.functions.getReserves().call()

    if is_reversed:
        peg_reserve = reserve0
        token_reserve = reserve1
    else:
        peg_reserve = reserve1
        token_reserve = reserve0

    if token_reserve and peg_reserve:
        if is_price_in_peg:
            # CALCULATE PRICE BY TOKEN PER PEG
            price = (Decimal(token_reserve) * 10 ** decimals) / (Decimal(peg_reserve) / ETHER)
        else:
            # CALCULATE PRICE BY PEG PER TOKEN
            price = (Decimal(peg_reserve) / ETHER) / (Decimal(token_reserve) * 10 ** decimals)

        return price, blockTimestampLast

    return Decimal('0'), -1

def get_bnb_price_of(token):

    cakeRouterV2 = Web3.toChecksumAddress('0x10ed43c718714eb63d5aa57b78b54704e256024e')
    cakeRouterV2 = web3.eth.contract(address=cakeRouterV2, abi=ABI).functions.factory().call()

    tokenChecksum = web3.toChecksumAddress(token)
    pair = web3.eth.contract(address=cakeRouterV2, abi=ABI).functions.getPair(tokenChecksum, WBNB).call()
    token_info = web3.eth.contract(tokenChecksum, abi=ERC20_ABI)
    token_symbol = token_info.functions.symbol().call()
    # print("decimals: ", token_decimals, "\n")
    pair_contract = web3.eth.contract(address=pair, abi=ABI)
    # print(pair_contract.functions.token0().call())
    is_reversed = pair_contract.functions.token0().call() == tokenChecksum
    decimals = web3.eth.contract(address=tokenChecksum, abi=ABI).functions.decimals().call()
    is_price_in_peg = True
    price, block = get_price(tokenChecksum, decimals, pair_contract, is_reversed, is_price_in_peg)
    return price, block, token_symbol

class PriceScraper:
    def __init__(self, web3, providerIndex):
        _instance = None
        self.web3 = web3
        self.mutex = Lock()
        self.count = 0
        self.providerIndex = providerIndex
        self.usdt_price = Decimal(0.0)
        self.block_timestamp_usdt = 0


    def get_price_json_of(self, address):
        block_timestamp = -1
        count = 0
        while count < len(providers):
            try:
                token_price, block_timestamp, symbol = get_bnb_price_of(address)
                providersSuccess[self.providerIndex] += 1
                break
            except:

                self.providerIndex += 1
                if self.providerIndex >= len(providers):
                    self.providerIndex = 0
                self.web3 = Web3(Web3.HTTPProvider(providers[self.providerIndex]))
                self.web3.middleware_onion.inject(geth_poa_middleware,
                                             layer=0)  # Again, this is needed for Binance, not Ethirium
                count += 1
        if block_timestamp == -1:
            print("failed", address)
            self.count += 1
            return
        block_timestamp = min(block_timestamp, self.block_timestamp_usdt)

        # print('1 Tarality =', token_price / self.usdt_price, 'USDT at block time', block_timestamp)
        block_iso8601_timestamp = get_iso8601_from_timestamp(block_timestamp)
        while True:
            priceData.append({
                'block': {
                    'height': -1,
                    'timestamp': {
                        'iso8601': block_iso8601_timestamp
                    }
                },
                'baseCurrency': {
                    'symbol': symbol,
                    'address': address
                },
                'quoteCurrency': {
                    'symbol': 'USDT',
                    'address': '0x55d398326f99059fF775485246999027B3197955'
                },
                'quotePrice': float(token_price / self.usdt_price)
            })
            self.count += 1
            # print(priceData)
            # print(get_iso8601_from_timestamp(block_timestamp))
            return

    def get_prices_of(self, addresses):
        for address in addresses:
            self.get_price_json_of(address)

if __name__ == '__main__':
    path = os.getcwd()
    providerIndex = int(sys.argv[3])
    print(path)
    web3 = Web3(Web3.HTTPProvider(providers[providerIndex]))
    web3.middleware_onion.inject(geth_poa_middleware, layer=0)  # Again, this is needed for Binance, not Ethirium
    print(f"Name of the script      : ", sys.argv[0])
    print(f"Arguments of the script : ", sys.argv[1])
    scraper = PriceScraper(web3, providerIndex=providerIndex)
    count = 0
    while count < 2*len(providers):
        try:
            scraper.usdt_price, scraper.block_timestamp_usdt, symbol = get_bnb_price_of('0x55d398326f99059fF775485246999027B3197955')
            providersSuccess[providerIndex] += 1
            break
        except Exception:

            providerIndex += 1
            if providerIndex >= len(providers):
                providerIndex = 0
            web3 = Web3(Web3.HTTPProvider(providers[providerIndex]))
            web3.middleware_onion.inject(geth_poa_middleware, layer=0)  # Again, this is needed for Binance, not Ethirium
            count += 1
    print('1 USDT =', scraper.usdt_price, 'BNB at block time', scraper.block_timestamp_usdt)
    usdt_iso8601_timestamp = get_iso8601_from_timestamp(scraper.block_timestamp_usdt)
    print(usdt_iso8601_timestamp)
    addresses = []
    file1 = open(sys.argv[1] + sys.argv[3], 'r')
    Lines = file1.readlines()
    address = '0x0fc812b96de7e910878039121938f6d5471b73dc'
    i = 0
    for line in Lines:
        address = line.strip()
        addresses.append(address)
        if i == 1:
            t = Thread(target=scraper.get_prices_of, args=[addresses])
            t.start()
            addresses = []
            i = -1

        i += 1
    t = Thread(target=scraper.get_prices_of, args=[addresses])
    t.start()
    while len(Lines) != scraper.count:
        time.sleep(1)
        print(len(priceData))
        print(scraper.count)
        print(len(Lines))
        print(providersSuccess)
        sys.stdout.flush()
    jsonContent ={
        'data': {
            'ethereum': {
                'dexTrades': priceData
            }
        }
    }
    if os.path.isfile(sys.argv[2]):
        os.remove(sys.argv[2])
    # os.remove(sys.argv[2])
    f = open(sys.argv[2], "w")
    jsonString = str(str(jsonContent).encode(sys.stdout.encoding, errors='replace')).replace("b\"", "").replace("}}}\"", "}}}").replace("'", "\"")
    f.write(jsonString)
    f.close()
    print(jsonString)
    sys.exit(0)

