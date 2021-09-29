package zcnsc

import (
	"0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/tokenpool"
	"0chain.net/core/common"
	"encoding/json"
)

// NewAuthorizerInfo To review: tokenLock init values
// pk = authorizer node public key
// authId = authorizer node public id = Client ID
func NewAuthorizerInfo(pk string, authId string, url string) *AuthorizerInfo {
	return &AuthorizerInfo{
		ID:        authId,
		PublicKey: pk,
		URL:       url,
		Staking: &tokenpool.ZcnLockingPool{
			ZcnPool: tokenpool.ZcnPool{
				TokenPool: tokenpool.TokenPool{
					ID:      "", // must be filled when DigPool is invoked. Usually this is a trx.Hash
					Balance: 0,  // filled when we dig pool
				},
			},
			TokenLockInterface: TokenLock{
				StartTime: 0,
				Duration:  0,
				Owner:     authId,
			},
		},
	}
}

func (an *AuthorizerInfo) Encode() []byte {
	bytes, _ := json.Marshal(an)
	return bytes
}

func (an *AuthorizerInfo) Decode(input []byte) error {
	tokenlock := &TokenLock{}

	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(input, &objMap)
	if err != nil {
		return err
	}

	id, ok := objMap["id"]
	if ok {
		var idStr *string
		err = json.Unmarshal(*id, &idStr)
		if err != nil {
			return err
		}
		an.ID = *idStr
	}

	pk, ok := objMap["public_key"]
	if ok {
		var pkStr *string
		err = json.Unmarshal(*pk, &pkStr)
		if err != nil {
			return err
		}
		an.PublicKey = *pkStr
	}

	url, ok := objMap["url"]
	if ok {
		var urlStr *string
		err = json.Unmarshal(*url, &urlStr)
		if err != nil {
			return err
		}
		an.URL = *urlStr
	}

	if an.Staking == nil {
		an.Staking = &tokenpool.ZcnLockingPool{
			ZcnPool: tokenpool.ZcnPool{
				TokenPool: tokenpool.TokenPool{},
			},
		}
	}

	staking, ok := objMap["staking"]
	if ok {
		err = an.Staking.Decode(*staking, tokenlock)
		if err != nil {
			return err
		}
	}
	return nil
}

func (an *AuthorizerInfo) Save(balances state.StateContextI) (err error) {
	_, err = balances.InsertTrieNode(ADDRESS+"auth_node"+an.ID, an)
	if err != nil {
		return common.NewError("save_auth_node_failed", "saving authorizer node: "+err.Error())
	}
	return nil
}
