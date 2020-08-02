# Blockchain

Blockchain implemented in Go which replicates the Bitcoin Blockchain. It is an open, distributed ledger that can record transactions efficiently by undergoing proof of work and consensus algorithm.

## Instructions to run

1. Clone this repo
```
git clone https://github.com/Sadham-Hussian/Blockchain.git
```
2. Set environment variable

    - In windows
        ```
        $ set NODE_ID:3000
        ```

    - In Linux
        ```
        $ export NODE_ID:3000
        ```

3. Create a wallet

    ```
    $ go run main.go createwallet
    ```
4. Create Blockchain and Genesis Block
    ```
    $ go run main.go CreateBlockchain -address ADDRESS
    ```
    Replace the address created by createwallet for ADDRESS in the above command

5. Start other nodes within a localnetwork 

If you have trouble setting up, please feel free to open an issue at GitHub

## Interactive Commands

```
1. $ go run main.go CreateBlockchain -address ADDRESS
2. $ go run main.go getBalance -address ADDRESS
3. $ go run main.go send -from FROM -to TO -amount AMOUNT -mine
4. $ go run main.go print
5. $ go run main.go createwallet
6. $ go run main.go listaddresses
7. $ go run main.go startnode 
```
To start a miner node enter the following command
```
$ go run main.go startnode -miner ADDRESS
```
