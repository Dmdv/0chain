package storagesc

import (
	"0chain.net/chaincore/chain/mocks"
	cstate "0chain.net/chaincore/chain/state"
	sci "0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/state"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
)

func TestPayBlobberBlockRewards(t *testing.T) {
	type parameters struct {
		blobberStakes         [][]float64
		blobberServiceCharge  []float64
		blockReward           float64
		qualifyingStake       float64
		blobberCapacityWeight float64
		blobberUsageWeight    float64
		minerWeight           float64
		sharderWeight         float64
	}

	type blockRewards struct {
		total                 state.Balance
		serviceChargeUsage    state.Balance
		serviceChargeCapacity state.Balance
		delegateUsage         []state.Balance
		delegateCapacity      []state.Balance
	}

	var blobberRewards = func(p parameters) []blockRewards {
		var totalQStake float64
		var rewards = []blockRewards{}
		var blobberStakes []float64
		for _, blobber := range p.blobberStakes {
			var stakes float64
			for _, dStake := range blobber {
				stakes += dStake
			}
			blobberStakes = append(blobberStakes, stakes)
			if stakes >= p.qualifyingStake {
				totalQStake += stakes
			}
		}

		var capacityRewardTotal = p.blockReward * p.blobberCapacityWeight /
			(p.blobberCapacityWeight + p.blobberUsageWeight + p.minerWeight + p.sharderWeight)
		var usageRewardTotal = p.blockReward * p.blobberUsageWeight /
			(p.blobberCapacityWeight + p.blobberUsageWeight + p.minerWeight + p.sharderWeight)

		for i, bStake := range blobberStakes {
			var reward blockRewards
			var capacityReward float64
			var usageReward float64
			var sc = p.blobberServiceCharge[i]
			if bStake >= p.qualifyingStake {
				capacityReward = capacityRewardTotal * bStake / totalQStake
				usageReward = usageRewardTotal * bStake / totalQStake

				reward.total = zcnToBalance(capacityReward + usageReward)
				reward.serviceChargeCapacity = zcnToBalance(capacityReward * sc)
				reward.serviceChargeUsage = zcnToBalance(usageReward * sc)
			}
			for _, dStake := range p.blobberStakes[i] {
				if bStake < p.qualifyingStake {
					reward.delegateUsage = append(reward.delegateUsage, 0)
					reward.delegateCapacity = append(reward.delegateCapacity, 0)
				} else {
					reward.delegateUsage = append(reward.delegateUsage, zcnToBalance(usageReward*(1-sc)*dStake/bStake))
					reward.delegateCapacity = append(reward.delegateCapacity, zcnToBalance(capacityReward*(1-sc)*dStake/bStake))
				}
			}
			rewards = append(rewards, reward)
		}
		return rewards
	}

	type args struct {
		ssc      *StorageSmartContract
		input    []byte
		balances cstate.StateContextI
	}

	var setExpectations = func(t *testing.T, p parameters) (*StorageSmartContract, cstate.StateContextI) {
		var balances = &mocks.StateContextI{}
		var blobbers = &StorageNodes{
			Nodes: []*StorageNode{},
		}
		var ssc = &StorageSmartContract{
			SmartContract: sci.NewSC(ADDRESS),
		}

		require.EqualValues(t, len(p.blobberStakes), len(p.blobberServiceCharge))
		rewards := blobberRewards(p)
		require.EqualValues(t, len(p.blobberStakes), len(rewards))

		var sPools []stakePool
		for i, reward := range rewards {
			require.EqualValues(t, len(reward.delegateCapacity), len(p.blobberStakes[i]))
			require.EqualValues(t, len(reward.delegateUsage), len(p.blobberStakes[i]))

			id := "bob " + strconv.Itoa(i)
			var sPool = stakePool{
				Pools: make(map[string]*delegatePool),
			}
			sPool.Settings.ServiceCharge = p.blobberServiceCharge[i]
			sPool.Settings.DelegateWallet = id
			sPools = append(sPools, sPool)
			require.True(t, blobbers.Nodes.add(&StorageNode{
				ID:                id,
				StakePoolSettings: stakePoolSettings{},
			}))

			for j := 0; j < len(reward.delegateUsage); j++ {
				require.EqualValues(t, len(reward.delegateUsage), len(reward.delegateCapacity))
				did := "delroy " + strconv.Itoa(j) + " " + id
				var dPool = &delegatePool{}
				dPool.ID = did
				dPool.Balance = zcnToBalance(p.blobberStakes[i][j])
				dPool.DelegateID = did
				sPool.Pools["paula "+did] = dPool
				if reward.delegateUsage[j] > 0 {
					balances.On("AddMint", &state.Mint{
						Minter: ADDRESS, ToClientID: did, Amount: reward.delegateUsage[j],
					}).Return(nil)
				}
				if reward.delegateCapacity[j] > 0 {
					balances.On("AddMint", &state.Mint{
						Minter: ADDRESS, ToClientID: did, Amount: reward.delegateCapacity[j],
					}).Return(nil)
				}
			}
			if reward.serviceChargeCapacity > 0 {
				balances.On("AddMint", &state.Mint{
					Minter: ADDRESS, ToClientID: id, Amount: reward.serviceChargeCapacity,
				}).Return(nil)
			}
			if reward.serviceChargeUsage > 0 {
				balances.On("AddMint", &state.Mint{
					Minter: ADDRESS, ToClientID: id, Amount: reward.serviceChargeUsage,
				}).Return(nil)
			}
			balances.On("GetTrieNode", stakePoolKey(ssc.ID, id)).Return(&sPool, nil).Once()
		}
		balances.On("GetTrieNode", ALL_BLOBBERS_KEY).Return(blobbers, nil).Once()
		var conf = &scConfig{
			BlockReward: &blockReward{
				BlockReward:           zcnToBalance(p.blockReward),
				QualifyingStake:       zcnToBalance(p.qualifyingStake),
				SharderWeight:         p.sharderWeight,
				MinerWeight:           p.minerWeight,
				BlobberCapacityWeight: p.blobberCapacityWeight,
				BlobberUsageWeight:    p.blobberUsageWeight,
			},
		}
		balances.On("GetTrieNode", scConfigKey(ssc.ID)).Return(conf, nil).Once()

		for i, sPool := range sPools {
			i := i
			sPool := sPool
			if rewards[i].total == 0 {
				continue
			}
			balances.On(
				"InsertTrieNode",
				stakePoolKey(ssc.ID, sPool.Settings.DelegateWallet),
				mock.MatchedBy(func(sp *stakePool) bool {
					for key, dPool := range sp.Pools {
						var wSplit = strings.Split(key, " ")
						dIndex, err := strconv.Atoi(wSplit[2])
						require.NoError(t, err)
						value, ok := sPools[i].Pools[key]

						if !ok ||
							value.DelegateID != dPool.DelegateID ||
							dPool.Rewards != rewards[i].delegateUsage[dIndex]+rewards[i].delegateCapacity[dIndex] {
							return false
						}
					}
					return sp.Rewards.Charge == rewards[i].serviceChargeCapacity+rewards[i].serviceChargeUsage &&
						sp.Rewards.Blobber == rewards[i].total &&
						sp.Settings.DelegateWallet == sPool.Settings.DelegateWallet

				}),
			).Return("", nil).Once()
		}
		conf.Minted += zcnToBalance(p.blockReward)
		balances.On("InsertTrieNode", scConfigKey(ssc.ID), conf).Return("", nil).Once()
		return ssc, balances
	}

	var parametersToArgs = func(t *testing.T, p parameters) args {
		ssc, balances := setExpectations(t, p)
		return args{
			ssc:      ssc,
			balances: balances,
		}
	}

	type want struct {
		error    bool
		errorMsg string
	}

	tests := []struct {
		name       string
		parameters parameters
		want       want
	}{
		{
			name: "ok",
			parameters: parameters{
				blobberStakes:         [][]float64{{10}, {5}},
				blobberServiceCharge:  []float64{0.1, 0.2},
				blockReward:           100,
				qualifyingStake:       10.0,
				blobberCapacityWeight: 5,
				blobberUsageWeight:    10,
				minerWeight:           2,
				sharderWeight:         3,
			},
		},

		{
			name: "ok",
			parameters: parameters{
				blobberStakes:         [][]float64{{10}, {10}},
				blobberServiceCharge:  []float64{0.1, 0.1},
				blockReward:           300,
				qualifyingStake:       10.0,
				blobberCapacityWeight: 5,
				blobberUsageWeight:    10,
				minerWeight:           2,
				sharderWeight:         3,
			},
		},
		{
			name: "ok",
			parameters: parameters{
				blobberStakes:         [][]float64{{10}, {5, 5, 10}, {2, 2}},
				blobberServiceCharge:  []float64{0.1, 0.2, 0.3},
				blockReward:           100,
				qualifyingStake:       10.0,
				blobberCapacityWeight: 5,
				blobberUsageWeight:    10,
				minerWeight:           2,
				sharderWeight:         3,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			args := parametersToArgs(t, tt.parameters)
			err := args.ssc.payBlobberBlockRewards(args.balances)

			require.EqualValues(t, tt.want.error, err != nil)
			if err != nil {
				require.EqualValues(t, tt.want.errorMsg, err.Error())
				return
			}
			require.True(t, mock.AssertExpectationsForObjects(t, args.balances))
		})
	}
}
