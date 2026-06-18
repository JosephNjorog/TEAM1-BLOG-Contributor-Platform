package avalanche

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// extractAmount pulls the uint256 amount argument back out of an ABI-
// encoded transfer(address,uint256) call: 4 byte selector + 32 byte
// address + 32 byte amount.
func extractAmount(t *testing.T, data []byte) *big.Int {
	t.Helper()
	if len(data) != 4+32+32 {
		t.Fatalf("unexpected encoded call length %d, want %d", len(data), 4+32+32)
	}
	return new(big.Int).SetBytes(data[4+32:])
}

func TestEncodeUSDCTransfer_ScalesToSixDecimals(t *testing.T) {
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	cases := []struct {
		amountUSD float64
		wantUnits int64
	}{
		{100.00, 100_000_000}, // the platform's fixed per-article rate
		{0.01, 10_000},        // one cent shouldn't underflow to zero
		{99.99, 99_990_000},   // a value that's awkward in binary float
		{1234.56, 1_234_560_000},
	}

	for _, c := range cases {
		data := encodeUSDCTransfer(to, c.amountUSD)
		got := extractAmount(t, data)
		want := big.NewInt(c.wantUnits)
		if got.Cmp(want) != 0 {
			t.Errorf("encodeUSDCTransfer(%.2f) amount = %s, want %s", c.amountUSD, got.String(), want.String())
		}
	}
}

func TestEncodeUSDCTransfer_IncludesRecipientAddress(t *testing.T) {
	to := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	data := encodeUSDCTransfer(to, 100.00)

	addressBytes := data[4 : 4+32]
	got := common.BytesToAddress(addressBytes)
	if got != to {
		t.Errorf("encoded recipient = %s, want %s", got.Hex(), to.Hex())
	}
}
