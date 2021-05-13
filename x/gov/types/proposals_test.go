package types_test

import (
	"fmt"
	govtypes "github.com/GeoDB-Limited/odin-core/x/gov/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProposalStatus_Format(t *testing.T) {
	statusDepositPeriod, _ := govtypes.ProposalStatusFromString("PROPOSAL_STATUS_DEPOSIT_PERIOD")
	tests := []struct {
		pt                   govtypes.ProposalStatus
		sprintFArgs          string
		expectedStringOutput string
	}{
		{statusDepositPeriod, "%s", "PROPOSAL_STATUS_DEPOSIT_PERIOD"},
		{statusDepositPeriod, "%v", "1"},
	}
	for _, tt := range tests {
		got := fmt.Sprintf(tt.sprintFArgs, tt.pt)
		require.Equal(t, tt.expectedStringOutput, got)
	}
}
