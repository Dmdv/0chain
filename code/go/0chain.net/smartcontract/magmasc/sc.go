package magmasc

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rcrowley/go-metrics"

	chainstate "0chain.net/chaincore/chain/state"
	sci "0chain.net/chaincore/smartcontractinterface"
	"0chain.net/chaincore/transaction"
	"0chain.net/core/common"
)

const (
	// Address is a SHA3-256 hex encoded hash of "magma" string. Represents address of MagmaSmartContract.
	Address = "11f8411db41e34cea7c100f19faff32da8f3cd5a80635731cec06f32d08089be"

	Name = "magma"
)

// MagmaSmartContract represents smartcontractinterface.SmartContractInterface implementation allows interacting with Magma.
type MagmaSmartContract struct {
	*sci.SmartContract
}

var (
	// Ensure MagmaSmartContract implements smartcontractinterface.SmartContractInterface.
	_ sci.SmartContractInterface = (*MagmaSmartContract)(nil)
)

// NewMagmaSmartContract creates smartcontractinterface.SmartContractInterface implemented by MagmaSmartContract with configured
// RestHandlers and SmartContractExecutionStats.
func NewMagmaSmartContract() sci.SmartContractInterface {
	var sscCopy = &MagmaSmartContract{
		SmartContract: sci.NewSC(Address),
	}
	sscCopy.setSC(sscCopy.SmartContract)
	return sscCopy
}

// These constants represents SmartContractExecutionStats keys, used to identify smart contract functions.
const (
	// registerConsumer represents name for Consumer's registration MagmaSmartContract function.
	registerConsumer = "register_consumer"

	// acceptTerms represents the name of MagmaSmartContract function.
	// When function is called it means that Consumer accepted Provider Terms.
	acceptTerms = "accept_terms"

	// registerProvider represents name for Provider's registration MagmaSmartContract function.
	registerProvider = "register_provider"
)

// setSC sets provided smartcontractinterface.SmartContract to corresponding MagmaSmartContract field
// and configures RestHandlers and SmartContractExecutionStats.
func (msc *MagmaSmartContract) setSC(sc *sci.SmartContract) {
	msc.SmartContract = sc

	// consumer
	msc.SmartContractExecutionStats[registerConsumer] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", msc.ID, registerConsumer), nil)
	msc.SmartContractExecutionStats[acceptTerms] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", msc.ID, acceptTerms), nil)
	msc.RestHandlers["/getAllConsumers"] = msc.getAllConsumers

	// provider
	msc.SmartContractExecutionStats[registerProvider] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", msc.ID, registerProvider), nil)
	msc.RestHandlers["/getProviderTerms"] = msc.getProviderTerms
	msc.SmartContractExecutionStats[acceptTerms] = metrics.GetOrRegisterTimer(fmt.Sprintf("sc:%v:func:%v", msc.ID, acceptTerms), nil)
}

// GetName implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetName() string {
	return Name
}

// GetAddress implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetAddress() string {
	return Address
}

// GetRestPoints implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetRestPoints() map[string]sci.SmartContractRestHandler {
	return msc.RestHandlers
}

// Execute implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) Execute(txn *transaction.Transaction,
	funcName string, input []byte, balances chainstate.StateContextI) (string, error) {

	switch funcName {
	// consumer
	case registerConsumer:
		return msc.registerConsumer(txn, balances)
	case acceptTerms:
		return msc.acceptTerms(txn, input, balances)

	// provider
	case registerProvider:
		return msc.registerProvider(txn, input, balances)
	}

	return "", common.NewError("invalid_function_name", "function with provided name is not supported")
}

// GetHandlerStats implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetHandlerStats(ctx context.Context, params url.Values) (interface{}, error) {
	return msc.SmartContract.HandlerStats(ctx, params)
}

// GetExecutionStats implements smartcontractinterface.SmartContractInterface.
func (msc *MagmaSmartContract) GetExecutionStats() map[string]interface{} {
	return msc.SmartContractExecutionStats
}
