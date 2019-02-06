package pbmoney

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.pedge.io/pb/go/pb/money"
)

func TestMath(t *testing.T) {
	money, err := pbmoney.NewMoneySimpleUSD(12, 0).
		Plus(pbmoney.NewMoneySimpleUSD(16, 0)).
		Minus(pbmoney.NewMoneySimpleUSD(8, 0)).
		TimesInt(2).
		DivInt(4).
		TimesFloat(2.0).
		Result()
	require.NoError(t, err)
	require.Equal(t, pbmoney.NewMoneySimpleUSD(20, 0), money)
	money, err = pbmoney.NewMoneySimpleUSD(12, 0).
		Plus(pbmoney.NewMoneySimpleUSD(16, 0)).
		Minus(pbmoney.NewMoneySimpleUSD(8, 0)).
		TimesInt(-2).
		DivInt(4).
		TimesFloat(2.0).
		PlusInt(500000).
		Result()
	require.NoError(t, err)
	require.Equal(t, pbmoney.NewMoneySimpleUSD(-19, 50), money)
}

func TestToFromGoogleMoney(t *testing.T) {
	testToFromGoogleMoney(t, pbmoney.NewMoneySimpleUSD(123456, 78))
	testToFromGoogleMoney(t, pbmoney.NewMoneySimpleUSD(-123456, 78))
}

func testToFromGoogleMoney(t *testing.T, money *pbmoney.Money) {
	money2, err := pbmoney.GoogleMoneyToMoney(money.ToGoogleMoney())
	fmt.Println(money.SimpleString())
	fmt.Println(money2.SimpleString())
	require.NoError(t, err)
	require.Equal(t, money, money2)
}
