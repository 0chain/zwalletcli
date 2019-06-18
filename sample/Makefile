SAMPLE_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

SCHEME=ed25519

include $(SAMPLE_DIR)/scheme/$(SCHEME).mk

ZWALLET=$(SAMPLE_DIR)/zwallet

CONFIG_DIR=$(SAMPLE_DIR)/config
CONFIG:=--configDir $(CONFIG_DIR) --config $(cluster)
ZWALLET_CMD:=$(ZWALLET) $(CONFIG)

FROM_WALLET:=--wallet $(from)
TO_WALLET:=--wallet $(to)

FAUCET:=faucet --methodName pour --input "Generous PayDay"

init: clean
	@echo ">>> Creating New Wallets"
	$(ZWALLET_CMD) getbalance $(FROM_WALLET)
	$(ZWALLET_CMD) getbalance $(TO_WALLET)

TO_CLIENT_ID=$(shell jq -r '.client_id' $(CONFIG_DIR)/$(to))
FROM_CLIENT_ID=$(shell jq -r '.client_id' $(CONFIG_DIR)/$(from))

transfer:
	@echo ">>> Balance before transfer"
	$(MAKE) getbalance
	@echo ""
	@echo ">>> Starting transfer"
	$(ZWALLET_CMD) $(FROM_WALLET) send --desc "Give loan - Make Happy" --toclientID $(TO_CLIENT_ID) --token 0.75
	@echo ""
	@echo ">>> Balance after transfer"
	$(MAKE) getbalance

repay:
	@echo ">>> Balance before transfer"
	$(MAKE) getbalance
	@echo ""
	@echo ">>> Starting transfer"
	$(ZWALLET_CMD) $(TO_WALLET) send --desc "Return loan - Feel sad" --toclientID $(FROM_CLIENT_ID) --token 0.50
	@echo ""
	@echo ">>> Balance after transfer"
	$(MAKE) getbalance

getbalance:
	$(ZWALLET_CMD) getbalance $(FROM_WALLET)
	$(ZWALLET_CMD) getbalance $(TO_WALLET)

FAUCET_LOOP?=10

faucet-loop:
	$(MAKE) getbalance
	for i in `seq $(FAUCET_LOOP)`; do $(MAKE) faucet; done
	$(MAKE) getbalance

faucet:
	@echo ">>> Add money receiver=$(from)"
	$(ZWALLET_CMD) $(FAUCET) $(FROM_WALLET)
	@echo ">>> Add money receiver=$(to)"
	$(ZWALLET_CMD) $(FAUCET) $(TO_WALLET)

clean:
	cd $(CONFIG_DIR);  rm -rf $(from) $(to)

setup:
	brew install jq || brew upgrade jq || true
	@mkdir -p $(ZCN_DIR) $(UPLOAD_DIR) $(DOWNLOAD_DIR) || true
	cp $(CLUSTERS)/*.yml $(ZCN_DIR)

help:
	@echo "Environment:"
	@echo "\tScheme            : $(SCHEME)"
	@echo "\tTo wallet         : $(CONFIG_DIR)/$(to)"
	@echo "\tFrom wallet       : $(CONFIG_DIR)/$(from)"
	@echo "\tCluster Config    : $(CONFIG_DIR)/$(cluster)"
	@echo "\tFaucet Loop Count : $(FAUCET_LOOP)"
	@echo
	@echo "Supported Commands: "
	@echo "Initialize Wallets"
	@echo "\tmake init           : Initialize the wallets"
	@echo ""
	@echo "Fund wallets"
	@echo "\tmake faucet         : Get initial pour of single token"
	@echo "\tmake faucet-loop    : Run faucet operation in a loop $(FAUCET_LOOP) times"
	@echo ""
	@echo "Transfer funds"
	@echo "\tmake transfer       : Send tokens 'from' client -> 'to' client"
	@echo "\tmake repay          : Send tokens 'to' client -> 'from' client"
	@echo ""
	@echo "Check Balance"
	@echo "\tmake getbalance     : Getbalance on all clients"
	@echo ""
	@echo "Clean up"
	@echo "\tmake clean          : Clean the wallets files"

