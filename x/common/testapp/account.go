package testapp

import (
	"math/rand"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Account is a data structure to store key of test account.
type Account struct {
	PrivKey    crypto.PrivKey
	PubKey     crypto.PubKey
	Address    sdk.AccAddress
	ValAddress sdk.ValAddress
}

func createArbitraryAccount(r *rand.Rand) Account {
	privkeySeed := make([]byte, 12)
	r.Read(privkeySeed)
	privKey := secp256k1.GenPrivKeySecp256k1(privkeySeed)
	return Account{
		PrivKey:    privKey,
		PubKey:     privKey.PubKey(),
		Address:    sdk.AccAddress(privKey.PubKey().Address()),
		ValAddress: sdk.ValAddress(privKey.PubKey().Address()),
	}
}
