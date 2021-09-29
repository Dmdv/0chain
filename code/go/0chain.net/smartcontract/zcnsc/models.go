package zcnsc

import (
	"0chain.net/chaincore/chain"
	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/config"
	"0chain.net/chaincore/state"
	"0chain.net/chaincore/transaction"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"0chain.net/core/common"
	"0chain.net/core/datastore"
	"0chain.net/core/encryption"
	"0chain.net/core/util"
)

var (
	AllAuthorizerKey = ADDRESS + encryption.Hash("all_authorizers")
)

type GlobalNode struct {
	ID                 string        `json:"id"`
	MinMintAmount      state.Balance `json:"min_mint_amount"`
	PercentAuthorizers float64       `json:"percent_authorizers"`
	MinBurnAmount      int64         `json:"min_burn_amount"`
	MinStakeAmount     int64         `json:"min_stake_amount"`
	BurnAddress        string        `json:"burn_address"`
	MinAuthorizers     int64         `json:"min_authorizers"`
}

func (gn *GlobalNode) GetKey() datastore.Key {
	return ADDRESS + gn.ID
}

func (gn *GlobalNode) GetHash() string {
	return util.ToHex(gn.GetHashBytes())
}

func (gn *GlobalNode) GetHashBytes() []byte {
	return encryption.RawHash(gn.Encode())
}

func (gn *GlobalNode) Encode() []byte {
	buff, _ := json.Marshal(gn)
	return buff
}

func (gn *GlobalNode) Decode(input []byte) error {
	err := json.Unmarshal(input, gn)
	return err
}

func (gn *GlobalNode) Save(balances cstate.StateContextI) (err error) {
	_, err = balances.InsertTrieNode(gn.GetKey(), gn)
	return
}

func GetGlobalSavedNode(balances cstate.StateContextI) (*GlobalNode, error) {
	gn := &GlobalNode{ID: ADDRESS}
	gv, err := balances.GetTrieNode(gn.GetKey())
	if err != nil {
		if err != util.ErrValueNotPresent {
			return nil, err
		} else {
			return gn, err
		}
	}
	if err := gn.Decode(gv.Encode()); err != nil {
		return nil, fmt.Errorf("%w: %s", common.ErrDecoding, err)
	}
	return gn, nil
}

// GetGlobalNode - returns global node
func GetGlobalNode(balances cstate.StateContextI) (*GlobalNode, error) {
	gn, err := GetGlobalSavedNode(balances)
	if err == nil {
		return gn, nil
	}

	if gn == nil {
		return nil, err
	}

	gn.MinMintAmount = state.Balance(config.SmartContractConfig.GetInt("smart_contracts.zcn.min_mint_amount"))
	gn.PercentAuthorizers = config.SmartContractConfig.GetFloat64("smart_contracts.zcn.percent_authorizers")
	gn.MinAuthorizers = config.SmartContractConfig.GetInt64("smart_contracts.zcn.min_authorizers")
	gn.MinBurnAmount = config.SmartContractConfig.GetInt64("smart_contracts.zcn.min_burn_amount")
	gn.MinStakeAmount = config.SmartContractConfig.GetInt64("smart_contracts.zcn.min_stake_amount")
	gn.BurnAddress = config.SmartContractConfig.GetString("smart_contracts.zcn.burn_address")

	return gn, nil
}

type MintPayload struct {
	EthereumTxnID     string                 `json:"ethereum_txn_id"`
	Amount            state.Balance          `json:"amount"`
	Nonce             int64                  `json:"nonce"`
	Signatures        []*AuthorizerSignature `json:"signatures"`
	ReceivingClientID string                 `json:"receiving_client_id"`
}

func (mp *MintPayload) Encode() []byte {
	buff, _ := json.Marshal(mp)
	return buff
}

func (mp *MintPayload) Decode(input []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(input, &objMap)
	if err != nil {
		return err
	}

	id, ok := objMap["ethereum_txn_id"]
	if ok {
		if id == nil {
			return errors.New("ethereum_txn_id is missing in the payload")
		}
		var value *string
		err = json.Unmarshal(*id, &value)
		if err != nil {
			return err
		}
		mp.EthereumTxnID = *value
	}

	id, ok = objMap["nonce"]
	if ok {
		if id == nil {
			return errors.New("nonce is missing in the payload")
		}
		var value *int64
		err = json.Unmarshal(*id, &value)
		if err != nil {
			return err
		}
		mp.Nonce = *value
	}

	id, ok = objMap["amount"]
	if ok {
		if id == nil {
			return errors.New("amount is missing in the payload")
		}
		var value *int64
		err = json.Unmarshal(*id, &value)
		if err != nil {
			return err
		}
		mp.Amount = state.Balance(*value)
	}

	id, ok = objMap["receiving_client_id"]
	if ok {
		if id == nil {
			return errors.New("receiving_client_id is missing in the payload")
		}
		var value *string
		err = json.Unmarshal(*id, &value)
		if err != nil {
			return err
		}
		mp.ReceivingClientID = *value
	}

	id, ok = objMap["signatures"]
	if ok {
		if id == nil {
			return errors.New("signatures is missing in the payload")
		}
		var sigs []*json.RawMessage
		err = json.Unmarshal(*id, &sigs)
		if err != nil {
			return err
		}

		for _, raw := range sigs {
			sig := &AuthorizerSignature{}
			err = json.Unmarshal(*raw, sig)
			if err != nil {
				return err
			}
			mp.Signatures = append(mp.Signatures, sig)
		}
	}

	return err
}

func (mp *MintPayload) GetStringToSign() string {
	return encryption.Hash(fmt.Sprintf("%v:%v:%v:%v", mp.EthereumTxnID, mp.Amount, mp.Nonce, mp.ReceivingClientID))
}

func (mp *MintPayload) verifySignatures(ans *AuthorizerNodes) (err error) {
	signatureScheme := chain.GetServerChain().GetSignatureScheme()
	toSign := mp.GetStringToSign()
	for _, v := range mp.Signatures {
		if v.ID == "" {
			return errors.New("authorizer ID is empty in a signature")
		}

		if ans.Nodes[v.ID] == nil {
			return errors.New(fmt.Sprintf("authorizer %s not found in authorizers", v.ID))
		}

		key := ans.Nodes[v.ID].PublicKey
		_ = signatureScheme.SetPublicKey(key)

		if key == "" {
			return errors.New("authorizer public key is empty")
		}

		ok, err := signatureScheme.Verify(v.Signature, toSign)
		if !ok || err != nil {
			return err
		}
	}

	return
}

type BurnPayloadResponse struct {
	TxnID           string `json:"0chain_txn_id"`
	Nonce           int64  `json:"nonce"`
	Amount          int64  `json:"amount"`
	EthereumAddress string `json:"ethereum_address"`
}

func (bp *BurnPayloadResponse) Encode() []byte {
	buff, _ := json.Marshal(bp)
	return buff
}

func (bp *BurnPayloadResponse) Decode(input []byte) error {
	err := json.Unmarshal(input, bp)
	return err
}

type BurnPayload struct {
	Nonce           int64  `json:"nonce"`
	EthereumAddress string `json:"ethereum_address"`
}

func (bp *BurnPayload) Encode() []byte {
	buff, _ := json.Marshal(bp)
	return buff
}

func (bp *BurnPayload) Decode(input []byte) error {
	err := json.Unmarshal(input, bp)
	return err
}

type UserNode struct {
	ID    string `json:"id"`
	Nonce int64  `json:"nonce"`
}

func (un *UserNode) GetKey(globalKey string) datastore.Key {
	return globalKey + un.ID
}

func (un *UserNode) GetHash() string {
	return util.ToHex(un.GetHashBytes())
}

func (un *UserNode) GetHashBytes() []byte {
	return encryption.RawHash(un.Encode())
}

func (un *UserNode) Encode() []byte {
	buff, _ := json.Marshal(un)
	return buff
}

func (un *UserNode) Decode(input []byte) error {
	err := json.Unmarshal(input, un)
	return err
}

func (un *UserNode) Save(balances cstate.StateContextI) (err error) {
	_, err = balances.InsertTrieNode(un.GetKey(ADDRESS), un)
	return
}

func GetUserNode(id string, balances cstate.StateContextI) (*UserNode, error) {
	un := &UserNode{ID: id}
	uv, err := balances.GetTrieNode(un.GetKey(ADDRESS))
	if err != nil {
		return un, err
	}
	if err := un.Decode(uv.Encode()); err != nil {
		return nil, fmt.Errorf("%w: %s", common.ErrDecoding, err)
	}
	return un, err
}

type TokenLock struct {
	StartTime common.Timestamp `json:"start_time"`
	Duration  time.Duration    `json:"duration"`
	Owner     datastore.Key    `json:"owner"`
}

func (tl TokenLock) IsLocked(entity interface{}) bool {
	txn, ok := entity.(*transaction.Transaction)
	if ok {
		return common.ToTime(txn.CreationDate).Sub(common.ToTime(tl.StartTime)) < tl.Duration
	}
	return true
}

func (tl TokenLock) LockStats(entity interface{}) []byte {
	txn, ok := entity.(*transaction.Transaction)
	if ok {
		p := &poolStat{
			StartTime: tl.StartTime,
			Duration:  tl.Duration,
			TimeLeft:  tl.Duration - common.ToTime(txn.CreationDate).Sub(common.ToTime(tl.StartTime)),
			Locked:    tl.IsLocked(txn),
		}
		return p.encode()
	}
	return nil
}

type poolStat struct {
	ID        datastore.Key    `json:"pool_id"`
	StartTime common.Timestamp `json:"start_time"`
	Duration  time.Duration    `json:"duration"`
	TimeLeft  time.Duration    `json:"time_left"`
	Locked    bool             `json:"locked"`
}

func (ps *poolStat) encode() []byte {
	buff, _ := json.Marshal(ps)
	return buff
}

func (ps *poolStat) decode(input []byte) error {
	err := json.Unmarshal(input, ps)
	return err
}
