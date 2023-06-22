package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"

	odin "github.com/ODIN-PROTOCOL/odin-core/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	legacy "github.com/cosmos/cosmos-sdk/types/bech32/legacybech32"
)

func Test_FromBech32(t *testing.T) {
	config := sdk.GetConfig()
	accountPrefix := odin.Bech32MainPrefix
	validatorPrefix := odin.Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	consensusPrefix := odin.Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	config.SetBech32PrefixForAccount(accountPrefix, accountPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForValidator(validatorPrefix, validatorPrefix+sdk.PrefixPublic)
	config.SetBech32PrefixForConsensusNode(consensusPrefix, consensusPrefix+sdk.PrefixPublic)
	bech32ConsPub := "odinvalconspub1addwnpepqge86lvslkpfk0rlz0ah9tat0vntx8yele36hhfpflehfehydlutkvdvhfm"
	mustConsPub, err := legacy.UnmarshalPubKey(sdk.PrefixConsensus+sdk.PrefixPublic, bech32ConsPub)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(mustConsPub.String())
	fmt.Println(mustConsPub.Type())
	bb := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, bb)
	encoder.Write(mustConsPub.Bytes())
	encoder.Close()
	fmt.Println(bb.String())
}
