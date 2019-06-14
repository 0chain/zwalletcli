PATH  := $(PATH):$(PWD)
SHELL := env PATH=$(PATH) /bin/bash

#SCHEME?=ed25519
SCHEME?=bls0chain

include $(SCHEME).mk

CONFIG:=--config $(cluster)
FROM_WALLET:=--wallet $(from)
TO_WALLET:=--wallet $(to)

ZCMD=zwallet $(CONFIG)

FAUCET:=faucet --methodName pour --input "Generous PayDay"

TO_CLIENT_ID:=$(shell jq -r '.client_id' $(HOME)/.zcn/$(to))
FROM_CLIENT_ID:=$(shell jq -r '.client_id' $(HOME)/.zcn/$(from))

init: clean
	@echo "Creating New Wallets"
	$(ZCMD) getbalance $(FROM_WALLET)
	$(ZCMD) getbalance $(TO_WALLET)

show:

send:
	$(ZCMD) $(FROM_WALLET) send --desc "Give loan - Make Happy" --toclientID $(TO_CLIENT_ID) --token 1.5

pay-chain: getbalance0 | send getbalance1 repay getbalance2

repay:
	$(ZCMD) $(TO_WALLET) send --desc "Return loan - Feel sad" --toclientID $(FROM_CLIENT_ID) --token 1.5

getbalance0 getbalance1 getbalance2 getbalance:
	$(ZCMD) getbalance $(FROM_WALLET)
	$(ZCMD) getbalance $(TO_WALLET)

getrich:
	for i in `seq 10`; do make faucet; done

faucet:
	@echo "Add money receiver=$(from)"
	$(ZCMD) $(FAUCET) $(FROM_WALLET)
	@echo "Add money receiver=$(to)"
	$(ZCMD) $(FAUCET) $(TO_WALLET)

clean:
	cd $(HOME)/.zcn;  rm $(from) $(to) || true


