package zcnsc

import (
	"0chain.net/chaincore/chain/state"
	"0chain.net/core/common"
	"0chain.net/core/datastore"
	"0chain.net/core/encryption"
	"0chain.net/core/logging"
	"0chain.net/core/util"
	"encoding/json"
	"go.uber.org/zap"

	//"github.com/pkg/errors"
	"fmt"
)

func FetchAuthorizers(balances state.StateContextI) (*AuthorizerNodes, error) {
	nodes := &AuthorizerNodes{
		Nodes: make(map[string]*AuthorizerInfo),
	}

	value, _ := balances.GetTrieNode(AllAuthorizerKey)
	if value == nil {
		return nodes, nil
	}

	bytes := value.Encode()
	err := nodes.Decode(bytes)
	logging.Logger.Info("getting authorizers node from MPT", zap.String("source", "fetchAuthorizers"), zap.String("nodes", string(bytes)))

	if err != nil {
		return nil, fmt.Errorf("%w: %s", common.ErrDecoding, err)
	}

	return nodes, nil
}

func (an *AuthorizerNodes) Exists(id datastore.Key) bool {
	return an.Nodes[id] != nil
}

func (an *AuthorizerNodes) DeleteAuthorizer(id string) (err error) {
	if an.Nodes[id] == nil {
		err = common.NewError("failed to delete authorizer", fmt.Sprintf("authorizer (%v) does not exist", id))
		return
	}
	delete(an.Nodes, id)
	return
}

func (an *AuthorizerNodes) AddAuthorizer(node *AuthorizerInfo) (err error) {
	if node == nil {
		err = common.NewError("failed to add authorizer", "authorizerNode is not initialized")
		return
	}

	if an.Nodes == nil {
		err = common.NewError("failed to add authorizer", "receiver Nodes is not initialized")
		return
	}

	if an.Nodes[node.ID] != nil {
		err = common.NewError("failed to add authorizer", fmt.Sprintf("authorizer (%v) already exists", node.ID))
		return
	}

	an.Nodes[node.ID] = node

	return
}

func (an *AuthorizerNodes) UpdateAuthorizer(node *AuthorizerInfo) (err error) {
	if an.Nodes[node.ID] == nil {
		err = common.NewError("failed to update authorizer", fmt.Sprintf("authorizer (%v) does not exist", node.ID))
		return
	}
	an.Nodes[node.ID] = node
	return
}

func (an *AuthorizerNodes) Decode(input []byte) error {
	if an.Nodes == nil {
		an.Nodes = make(map[string]*AuthorizerInfo)
	}

	var objMap map[string]json.RawMessage
	err := json.Unmarshal(input, &objMap)
	if err != nil {
		return err
	}

	nodeMap, ok := objMap[Nodes]
	if ok {
		var authorizerNodes map[string]json.RawMessage
		err := json.Unmarshal(nodeMap, &authorizerNodes)
		if err != nil {
			return err
		}

		for _, raw := range authorizerNodes {
			target := &AuthorizerInfo{}
			err := target.Decode(raw)
			if err != nil {
				return err
			}

			err = an.AddAuthorizer(target)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (an *AuthorizerNodes) Encode() []byte {
	buff, _ := json.Marshal(an)
	return buff
}

func (an *AuthorizerNodes) GetHash() string {
	return util.ToHex(an.GetHashBytes())
}

func (an *AuthorizerNodes) GetHashBytes() []byte {
	return encryption.RawHash(an.Encode())
}

func (an *AuthorizerNodes) Save(balances state.StateContextI) (err error) {
	_, err = balances.InsertTrieNode(AllAuthorizerKey, an)
	return
}
