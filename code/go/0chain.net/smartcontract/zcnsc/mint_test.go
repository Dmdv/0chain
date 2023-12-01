package zcnsc_test

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"0chain.net/chaincore/chain"
	cstate "0chain.net/chaincore/chain/state"
	"0chain.net/chaincore/state"
	"0chain.net/core/common"
	"0chain.net/core/config"
	"0chain.net/smartcontract/dbs/event"
	"0chain.net/smartcontract/partitions"
	"0chain.net/smartcontract/stakepool"
	"0chain.net/smartcontract/stakepool/spenum"
	. "0chain.net/smartcontract/zcnsc"
	"github.com/0chain/common/core/currency"
	"github.com/0chain/common/core/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	logging.Logger = zap.NewNop()

	config.SetupDefaultConfig()

	chainConfig := chain.NewConfigImpl(&chain.ConfigData{})
	chainConfig.FromViper() //nolint: errcheck

	config.Configuration().ChainConfig = chainConfig
}

func Test_MintPayload_Encode_Decode(t *testing.T) {
	ctx := MakeMockStateContext()
	expected, err := CreateMintPayload(ctx, defaultClient)
	require.NoError(t, err)
	actual := &MintPayload{}
	err = actual.Decode(expected.Encode())
	require.NoError(t, err)
	require.Equal(t, expected.Nonce, actual.Nonce)
	require.Equal(t, expected.Amount, actual.Amount)
	require.Equal(t, expected.EthereumTxnID, actual.EthereumTxnID)
	require.Equal(t, expected.ReceivingClientID, actual.ReceivingClientID)
	require.Equal(t, len(expected.Signatures), len(actual.Signatures))
	for i := range actual.Signatures {
		require.Equal(t, expected.Signatures[i].ID, actual.Signatures[i].ID)
		require.Equal(t, expected.Signatures[i].Signature, actual.Signatures[i].Signature)
	}
}

func Test_DifferentSenderAndReceiverMustFail(t *testing.T) {
	ctx := MakeMockStateContext()
	contract := CreateZCNSmartContract()

	payload, err := CreateMintPayload(ctx, defaultClient)
	require.NoError(t, err)

	transaction, err := CreateTransaction(defaultClient+"1", "mint", payload.Encode(), ctx)
	require.NoError(t, err)

	eventDb, err := event.NewInMemoryEventDb(config.DbAccess{}, config.DbSettings{
		Debug:                 true,
		PartitionChangePeriod: 1,
	})
	require.NoError(t, err)

	err = eventDb.Get().Model(&event.User{}).Create(&event.User{
		UserID:    transaction.ClientID,
		MintNonce: 0,
	}).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		err = eventDb.Drop()
		require.NoError(t, err)

		eventDb.Close()
	})

	ctx.SetEventDb(eventDb)

	_, err = contract.Mint(transaction, payload.Encode(), ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "transaction made from different account who made burn")
}

func Test_MaxFeeMint(t *testing.T) {
	type expect struct {
		sharedFee    currency.Coin
		remainAmount currency.Coin
	}

	tt := []struct {
		name   string
		maxFee currency.Coin
		expect expect
	}{
		{
			name:   "max fee not evenly distributed",
			maxFee: 10,
			expect: expect{
				sharedFee:    3,
				remainAmount: 197,
			},
		},
		{
			name:   "max fee evenly distributed",
			maxFee: 9,
			expect: expect{
				sharedFee:    3,
				remainAmount: 197,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := MakeMockStateContext()
			ctx.globalNode.ZCNSConfig.MaxFee = tc.maxFee
			contract := CreateZCNSmartContract()
			payload, err := CreateMintPayload(ctx, defaultClient)
			require.NoError(t, err)

			transaction, err := CreateTransaction(defaultClient, "mint", payload.Encode(), ctx)
			require.NoError(t, err)

			eventDb, err := event.NewInMemoryEventDb(config.DbAccess{}, config.DbSettings{
				Debug:                 true,
				PartitionChangePeriod: 1,
			})
			require.NoError(t, err)

			err = eventDb.Get().Model(&event.User{}).Create(&event.User{
				UserID:    transaction.ClientID,
				MintNonce: 0,
			}).Error
			require.NoError(t, err)

			t.Cleanup(func() {
				err = eventDb.Drop()
				require.NoError(t, err)

				eventDb.Close()
			})

			ctx.SetEventDb(eventDb)

			response, err := contract.Mint(transaction, payload.Encode(), ctx)
			require.NoError(t, err, "Testing authorizer: '%s'", defaultClient)
			require.NotNil(t, response)
			require.NotEmpty(t, response)

			mm := ctx.GetTransfers()
			require.Equal(t, 1, len(mm))

			auths := make([]string, 0, len(payload.Signatures))
			for _, sig := range payload.Signatures {
				auths = append(auths, sig.ID)
			}

			mintsMap := make(map[string]*state.Transfer, len(mm))
			for i, m := range mm {
				mintsMap[m.ToClientID] = mm[i]
			}

			// sort the signatures
			sortedSigs := make([]*AuthorizerSignature, len(payload.Signatures))
			copy(sortedSigs, payload.Signatures)
			sort.Slice(sortedSigs, func(i, j int) bool {
				return sortedSigs[i].ID < sortedSigs[j].ID
			})

			rd := rand.New(rand.NewSource(ctx.GetBlock().GetRoundRandomSeed()))
			sig := sortedSigs[rd.Intn(len(payload.Signatures))]

			stakePool := NewStakePool()
			err = ctx.GetTrieNode(stakepool.StakePoolKey(spenum.Authorizer, sig.ID), stakePool)
			require.NoError(t, err)
			require.Equal(t, tc.expect.sharedFee, stakePool.Reward)

			// assert transaction.ClientID has remaining amount
			tm, ok := mintsMap[transaction.ClientID]
			require.True(t, ok)
			require.Equal(t, tc.expect.remainAmount, tm.Amount)
		})
	}
}

func Test_EmptySignaturesShouldFail(t *testing.T) {
	ctx := MakeMockStateContext()
	contract := CreateZCNSmartContract()
	payload, err := CreateMintPayload(ctx, defaultClient)
	require.NoError(t, err)

	payload.Signatures = nil

	transaction, err := CreateTransaction(defaultClient, "mint", payload.Encode(), ctx)
	require.NoError(t, err)

	eventDb, err := event.NewInMemoryEventDb(config.DbAccess{}, config.DbSettings{
		Debug:                 true,
		PartitionChangePeriod: 1,
	})
	require.NoError(t, err)

	err = eventDb.Get().Model(&event.User{}).Create(&event.User{
		UserID:    transaction.ClientID,
		MintNonce: 0,
	}).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		err = eventDb.Drop()
		require.NoError(t, err)

		eventDb.Close()
	})

	ctx.SetEventDb(eventDb)

	_, err = contract.Mint(transaction, payload.Encode(), ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "signatures entry is missing in payload")
}

func Test_EmptyAuthorizersNonemptySignaturesShouldFail(t *testing.T) {
	ctx := MakeMockStateContextWithoutAutorizers()

	contract := CreateZCNSmartContract()
	payload, err := CreateMintPayload(ctx, defaultClient)
	require.NoError(t, err)

	// Add a few signatures.
	var signatures []*AuthorizerSignature
	for _, id := range []string{"sign1", "sign2", "sign3"} {
		signatures = append(signatures, &AuthorizerSignature{ID: id})
	}
	payload.Signatures = signatures

	transaction, err := CreateTransaction(defaultClient, "mint", payload.Encode(), ctx)
	require.NoError(t, err)

	eventDb, err := event.NewInMemoryEventDb(config.DbAccess{}, config.DbSettings{
		Debug:                 true,
		PartitionChangePeriod: 1,
	})
	require.NoError(t, err)

	err = eventDb.Get().Model(&event.User{}).Create(&event.User{
		UserID:    transaction.ClientID,
		MintNonce: 0,
	}).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		err = eventDb.Drop()
		require.NoError(t, err)

		eventDb.Close()
	})

	ctx.SetEventDb(eventDb)

	_, err = contract.Mint(transaction, payload.Encode(), ctx)
	require.Equal(t, common.NewError("failed to mint", "no authorizers found"), err)
}

func Test_MintPayloadNonceShouldBeRecordedByUserNode(t *testing.T) {
	ctx := MakeMockStateContext()

	tr := CreateDefaultTransactionToZcnsc()
	eventDb, err := event.NewInMemoryEventDb(config.DbAccess{}, config.DbSettings{
		Debug:                 true,
		PartitionChangePeriod: 1,
	})
	require.NoError(t, err)

	err = eventDb.Get().Model(&event.User{}).Create(&event.User{
		UserID:    tr.ClientID,
		MintNonce: 0,
	}).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		err = eventDb.Drop()
		require.NoError(t, err)

		eventDb.Close()
	})

	ctx.SetEventDb(eventDb)

	payload, err := CreateMintPayload(ctx, defaultClient)
	require.NoError(t, err)

	contract := CreateZCNSmartContract()

	payload.Nonce = 1

	resp, err := contract.Mint(tr, payload.Encode(), ctx)
	require.NoError(t, err)
	require.NotZero(t, resp)

	user, err := ctx.GetEventDB().GetUser(tr.ClientID)
	require.NoError(t, err)
	require.Equal(t, payload.Nonce, user.MintNonce)
}

func Test_CheckAuthorizerStakePoolDistributedRewards(t *testing.T) {
	ctx := MakeMockStateContext()

	tr := CreateDefaultTransactionToZcnsc()
	eventDb, err := event.NewInMemoryEventDb(config.DbAccess{}, config.DbSettings{
		Debug:                 true,
		PartitionChangePeriod: 1,
	})
	require.NoError(t, err)

	err = eventDb.Get().Model(&event.User{}).Create(&event.User{
		UserID:    tr.ClientID,
		MintNonce: 0,
	}).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		err = eventDb.Drop()
		require.NoError(t, err)

		eventDb.Close()
	})

	ctx.SetEventDb(eventDb)

	payload, err := CreateMintPayload(ctx, defaultClient)
	require.NoError(t, err)

	contract := CreateZCNSmartContract()

	payload.Nonce = 1

	gn, err := GetGlobalNode(ctx)
	require.NoError(t, err)

	gn.ZCNSConfig.MaxFee = 100
	err = gn.Save(ctx)
	require.NoError(t, err)

	rand.Seed(ctx.GetBlock().GetRoundRandomSeed())

	rewardBefore := 0

	resp, err := contract.Mint(tr, payload.Encode(), ctx)
	require.NoError(t, err)
	require.NotZero(t, resp)

	rewardAfter := 0

	for _, sig := range payload.Signatures {
		stakePool := NewStakePool()
		err = ctx.GetTrieNode(stakepool.StakePoolKey(spenum.Authorizer, sig.ID), stakePool)
		require.NoError(t, err)

		rewardAfter += int(stakePool.Reward)
	}

	require.Greater(t, rewardAfter, rewardBefore, "reward should be distributed rewardAfter : ? rewardBefore : ?", rewardAfter, rewardBefore)
}

func TestZCNSmartContractMintNonce(t *testing.T) {
	tt := []struct {
		name         string
		nonce        int64
		initFunc     func(state cstate.StateContextI)
		expectErrMsg string  // substring of expected error message
		expectNonces []int64 // nonces should be saved to the partitions
	}{
		{
			name:  "add nonce - ok",
			nonce: 1,
		},
		{
			name: "add new nonce - ok",
			initFunc: func(state cstate.StateContextI) {
				err := PartitionWZCNMintedNonceAdd(state, 1)
				require.NoError(t, err)
			},
			nonce:        2,
			expectErrMsg: "",
		},
		{
			name: "add duplicate nonce - fail",
			initFunc: func(state cstate.StateContextI) {
				err := PartitionWZCNMintedNonceAdd(state, 1)
				require.NoError(t, err)
			},
			nonce:        1,
			expectErrMsg: "has already been minted",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := MakeMockStateContext()
			ctx.globalNode.ZCNSConfig.MaxFee = 10
			contract := CreateZCNSmartContract()
			payload, err := CreateMintPayloadWithNonce(ctx, defaultClient, tc.nonce)
			require.NoError(t, err)

			transaction, err := CreateTransaction(defaultClient, "mint", payload.Encode(), ctx)
			require.NoError(t, err)

			eventDb, err := event.NewInMemoryEventDb(config.DbAccess{}, config.DbSettings{
				Debug:                 true,
				PartitionChangePeriod: 1,
			})
			require.NoError(t, err)

			err = eventDb.Get().Model(&event.User{}).Create(&event.User{
				UserID:    transaction.ClientID,
				MintNonce: 0,
			}).Error
			require.NoError(t, err)

			t.Cleanup(func() {
				err = eventDb.Drop()
				require.NoError(t, err)

				eventDb.Close()
			})

			ctx.SetEventDb(eventDb)

			if tc.initFunc != nil {
				tc.initFunc(ctx)
			}

			_, err = contract.Mint(transaction, payload.Encode(), ctx)
			if tc.expectErrMsg != "" {
				require.ErrorContains(t, err, tc.expectErrMsg)
			}
			if err != nil {
				return
			}

			mm := ctx.GetTransfers()
			require.Equal(t, 1, len(mm))

			// check that the nonce is saved to the partition by calling the Add and got
			// error of item already exists
			err = PartitionWZCNMintedNonceAdd(ctx, tc.nonce)
			require.True(t, partitions.ErrItemExist(err))
		})
	}
}
