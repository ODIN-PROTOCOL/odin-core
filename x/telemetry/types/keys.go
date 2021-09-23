package types

const (
	// ModuleName is the name of the module.
	ModuleName = "telemetry"
	// StoreKey to be used when creating the KVStore.
	StoreKey     = ModuleName
	RouterKey    = ModuleName
	QuerierRoute = ModuleName

	QueryTopBalances        = "top_balances"
	QueryExtendedValidators = "extended_validators"
	QueryAvgBlockSize       = "avg_block_size"
	QueryAvgBlockTime       = "avg_block_time"
	QueryAvgTxFee           = "avg_tx_fee"
	QueryTxVolume           = "tx_volume"
	QueryValidatorBlocks    = "validator_blocks"
	QueryTopValidators      = "top_validators"
	QueryBalances           = "balances"

	DenomTag  = "denom"
	StatusTag = "status"
)
