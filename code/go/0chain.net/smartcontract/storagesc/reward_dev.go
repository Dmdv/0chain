package storagesc

import (
	"0chain.net/core/common"
	"0chain.net/smartcontract/stakepool/spenum"
	"net/http"
	"strconv"
)

func (srh *StorageRestHandler) getAllChallenges(w http.ResponseWriter, r *http.Request) {
	// read all data from challenges table and return
	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}

	allocationID := r.URL.Query().Get("allocation_id")

	challenges, err := edb.GetAllChallengesByAllocationID(allocationID)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	common.Respond(w, r, challenges, nil)
}

func (srh *StorageRestHandler) getBlockRewards(w http.ResponseWriter, r *http.Request) {
	// read all data from block_rewards table and return
	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}

	startBlockNumber := r.URL.Query().Get("start_block_number")
	endBlockNumber := r.URL.Query().Get("end_block_number")

	result, err := edb.GetBlockRewards(startBlockNumber, endBlockNumber)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	common.Respond(w, r, result, nil)
}

func (srh *StorageRestHandler) getReadRewards(w http.ResponseWriter, r *http.Request) {
	// read all data from challenge_rewards table and return
	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}

	allocationID := r.URL.Query().Get("allocation_id")

	result, err := edb.GetAllocationReadRewards(allocationID)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	common.Respond(w, r, result, nil)
}

func (srh *StorageRestHandler) getTotalChallengeRewards(w http.ResponseWriter, r *http.Request) {
	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}

	allocationID := r.URL.Query().Get("allocation_id")

	totalBlobberRewards := map[string]int64{}
	totalValidatorRewards := map[string]int64{}

	challengeRewards, err := edb.GetAllocationChallengeRewards(allocationID)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	for i, j := range challengeRewards {
		if int(j.ProviderType) == spenum.ChallengePassReward.Int() {
			totalBlobberRewards[i] = j.Total
		} else {
			totalValidatorRewards[i] = j.Total
		}
	}

	result := map[string]interface{}{
		"blobber_rewards":   totalBlobberRewards,
		"validator_rewards": totalValidatorRewards,
	}

	common.Respond(w, r, result, nil)
}

func (srh *StorageRestHandler) getAllocationCancellationReward(w http.ResponseWriter, r *http.Request) {
	// read all data from allocation_cancellation_rewards table and return
	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}

	allocationID := r.URL.Query().Get("allocation_id")

	result, err := edb.GetAllocationCancellationRewards(allocationID)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	common.Respond(w, r, result, nil)
}

func (srh *StorageRestHandler) getAllocationChallengeRewards(w http.ResponseWriter, r *http.Request) {

	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}

	allocationID := r.URL.Query().Get("allocation_id")

	result, err := edb.GetAllocationChallengeRewards(allocationID)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	common.Respond(w, r, result, err)
}

func (srh *StorageRestHandler) getPassedChallengesForBlobberAllocation(w http.ResponseWriter, r *http.Request) {

	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}

	allocationID := r.URL.Query().Get("allocation_id")

	result, err := edb.GetPassedChallengesForBlobberAllocation(allocationID)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	common.Respond(w, r, result, err)
}

func (srh *StorageRestHandler) getChallengesCountBetweenBlocks(w http.ResponseWriter, r *http.Request) {
	edb := srh.GetQueryStateContext().GetEventDB()
	if edb == nil {
		common.Respond(w, r, nil, common.NewErrInternal("no db connection"))
		return
	}

	startString := r.URL.Query().Get("start")
	endString := r.URL.Query().Get("end")

	start, err := strconv.ParseInt(startString, 10, 64)
	if err != nil {
		common.Respond(w, r, nil, common.NewErrBadRequest("invalid start block number"))
		return
	}

	end, err := strconv.ParseInt(endString, 10, 64)
	if err != nil {
		common.Respond(w, r, nil, common.NewErrBadRequest("invalid end block number"))
		return
	}

	result, err := edb.GetChallengesCountBetweenBlocks(start, end)
	if err != nil {
		common.Respond(w, r, nil, err)
		return
	}

	common.Respond(w, r, result, err)
}
