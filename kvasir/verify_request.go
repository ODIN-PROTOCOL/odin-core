package kvasir

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type RequestVerification struct {
	ChainID   string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	Validator string `protobuf:"bytes,2,opt,name=validator,proto3" json:"validator,omitempty"`
	RequestID uint64 `protobuf:"varint,3,opt,name=request_id,json=requestId,proto3,casttype=RequestID" json:"request_id,omitempty"`
	Contract  string `protobuf:"varint,4,opt,name=contract,json=externalId,proto3,casttype=ExternalID" json:"contract,omitempty"`
}

func NewRequestVerification(
	chainID string,
	validator sdk.ValAddress,
	requestID uint64,
	contract string,
) RequestVerification {
	return RequestVerification{
		ChainID:   chainID,
		Validator: validator.String(),
		RequestID: requestID,
		Contract:  contract,
	}
}

func (msg RequestVerification) GetSignBytes() []byte {
	bz, _ := json.Marshal(msg)
	return sdk.MustSortJSON(bz)
}
