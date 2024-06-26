package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ODIN-PROTOCOL/odin-core/testing/testapp"
	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/keeper"
)

func TestSnapshotter(t *testing.T) {
	// setup source app
	srcApp, srcCtx, srcKeeper := testapp.CreateTestInput(true)

	// create snapshot
	_, err := srcApp.Commit()
	require.NoError(t, err)

	srcHashToCode, err := getMappingHashToCode(srcCtx, &srcKeeper)
	require.NoError(t, err)
	snapshotHeight := uint64(srcApp.LastBlockHeight())
	snapshot, err := srcApp.SnapshotManager().Create(snapshotHeight)
	require.NoError(t, err)
	assert.NotNil(t, snapshot)

	// restore snapshot
	destApp := testapp.SetupWithEmptyStore()
	destCtx := destApp.NewUncachedContext(false, tmproto.Header{})
	destKeeper := destApp.OracleKeeper
	require.NoError(t, destApp.SnapshotManager().Restore(*snapshot))
	for i := uint32(0); i < snapshot.Chunks; i++ {
		chunkBz, err := srcApp.SnapshotManager().LoadChunk(snapshot.Height, snapshot.Format, i)
		require.NoError(t, err)
		end, err := destApp.SnapshotManager().RestoreChunk(chunkBz)
		require.NoError(t, err)
		if end {
			break
		}
	}
	destHashToCode, err := getMappingHashToCode(destCtx, &destKeeper)
	require.NoError(t, err)

	// compare src and dest
	assert.Equal(
		t,
		srcHashToCode,
		destHashToCode,
	)
}

func getMappingHashToCode(ctx sdk.Context, keeper *keeper.Keeper) (map[string][]byte, error) {
	hashToCode := make(map[string][]byte)
	oracleScripts, err := keeper.GetAllOracleScripts(ctx)
	if err != nil {
		return nil, err
	}
	for _, oracleScript := range oracleScripts {
		hashToCode[oracleScript.Filename] = keeper.GetFile(oracleScript.Filename)
	}
	dataSources, err := keeper.GetAllDataSources(ctx)
	if err != nil {
		return hashToCode, err
	}
	for _, dataSource := range dataSources {
		hashToCode[dataSource.Filename] = keeper.GetFile(dataSource.Filename)
	}

	return hashToCode, nil
}
