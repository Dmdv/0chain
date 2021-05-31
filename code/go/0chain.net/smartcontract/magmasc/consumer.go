package magmasc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
	"0chain.net/core/util"
)

// registerConsumer represents registerConsumer MagmaSmartContract function and allows registering Consumer in blockchain.
//
// registerConsumer creates Consumer with Consumer.ID (equals to transaction client ID),
// adds it to all Consumers list, creates stakePool for new Consumer and saves results in provided state.StateContextI.
func (msc *MagmaSmartContract) registerConsumer(txn *transaction.Transaction, balances state.StateContextI) (string, error) {
	const errCode = "register_consumer"

	consumers, err := extractConsumers(balances)
	if err != nil {
		return "", common.NewErrorf(errCode, "retrieving all consumers from state failed with error: %v ", err)
	}

	var (
		consumer = Consumer{
			ID: txn.ClientID,
		}
	)
	if containsConsumer(msc.ID, consumer, consumers, balances) {
		return "", common.NewErrorf(errCode, "consumer with id=`%s` already exist", consumer.ID)
	}

	if err := createAndInsertConsumerStakePool(consumer.ID, msc.ID, balances); err != nil {
		return "", common.NewErrorf(errCode, "creating stake pool for consumer failed with err: %v", err)
	}

	// save the all consumers
	consumers.Nodes.add(&consumer)
	_, err = balances.InsertTrieNode(AllConsumersKey, consumers)
	if err != nil {
		return "", common.NewErrorf(errCode, "saving the all consumers failed with error: %v ", err)
	}

	// save the new consumer
	_, err = balances.InsertTrieNode(nodeKey(msc.ID, consumer.ID, consumerType), &consumer)
	if err != nil {
		return "", common.NewErrorf(errCode, "saving consumer failed with error: %v ", err)
	}

	return string(consumer.Encode()), nil
}

// extractConsumers extracts all consumers represented in JSON bytes stored in state.StateContextI with AllConsumersKey.
//
// extractConsumers returns err if state.StateContextI does not contain consumers or stored bytes have invalid format.
func extractConsumers(balances state.StateContextI) (*Consumers, error) {
	consumers := &Consumers{}
	consumerTN, err := balances.GetTrieNode(AllConsumersKey)
	if err != nil && err != util.ErrValueNotPresent {
		return nil, err
	}
	if err == util.ErrValueNotPresent || consumerTN == nil {
		return consumers, nil
	}

	if err := json.Unmarshal(consumerTN.Encode(), consumers); err != nil {
		return nil, fmt.Errorf("%w: %s", common.ErrDecoding, err)
	}

	return consumers, nil
}

// createAndInsertConsumerStakePool creates stakePool for Consumer and saves it in state.StateContextI.
//
// if stakePool for provided Consumer.ID already exist it returns ErrStakePoolExist. Also, createAndInsertConsumerStakePool
// returns err occurred while inserting new stakePool in state.StateContextI.
func createAndInsertConsumerStakePool(consumerID, scKey string, balances state.StateContextI) error {
	_, err := balances.GetTrieNode(stakePoolKey(scKey, consumerID))
	if err != util.ErrValueNotPresent {
		return ErrStakePoolExist
	}

	sp := newStakePool()
	sp.ID = consumerID

	_, err = balances.InsertTrieNode(stakePoolKey(scKey, consumerID), sp)
	if err != nil {
		return err
	}

	return nil
}

// getAllConsumers represents MagmaSmartContract handler. Returns all registered Consumer's nodes
// stores in provided state.StateContextI with AllConsumersKey.
func (msc *MagmaSmartContract) getAllConsumers(_ context.Context, _ url.Values, balances state.StateContextI) (interface{}, error) {
	consumers, err := extractConsumers(balances)
	if err != nil {
		return "", common.NewErrInternal("err while extracting all consumers from state", err.Error())
	}

	return consumers.Nodes, nil
}
