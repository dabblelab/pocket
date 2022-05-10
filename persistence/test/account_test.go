package test

import (
	"math/big"
	"testing"

	"github.com/pokt-network/pocket/persistence"
	"github.com/pokt-network/pocket/shared/crypto"
	"github.com/pokt-network/pocket/shared/types"
	typesGenesis "github.com/pokt-network/pocket/shared/types/genesis"
	"github.com/stretchr/testify/require"
)

// --- Account Tests ---

func TestSetAccountAmount(t *testing.T) {
	db := persistence.PostgresContext{
		Height: 0,
		DB:     *PostgresDB,
	}
	acc := NewTestAccount(t)

	err := db.SetAccountAmount(acc.Address, DefaultStake)
	require.NoError(t, err)

	am, err := db.GetAccountAmount(acc.Address)
	require.NoError(t, err)
	require.Equal(t, DefaultStake, am, "unexpected amount")

	db.SetAccountAmount(acc.Address, StakeToUpdate)
	require.NoError(t, err)

	am, err = db.GetAccountAmount(acc.Address)
	require.NoError(t, err)
	require.Equal(t, StakeToUpdate, am, "unexpected amount after second set")
}

func TestAddAccountAmount(t *testing.T) {
	db := persistence.PostgresContext{
		Height: 0,
		DB:     *PostgresDB,
	}
	acc := NewTestAccount(t)

	err := db.SetAccountAmount(acc.Address, DefaultStake)
	require.NoError(t, err)

	amountToAddBig := big.NewInt(100)
	err = db.AddAccountAmount(acc.Address, types.BigIntToString(amountToAddBig))
	require.NoError(t, err)

	am, err := db.GetAccountAmount(acc.Address)
	require.NoError(t, err)

	resultBig := (&big.Int{}).Add(DefaultStakeBig, amountToAddBig)
	expectedResult := types.BigIntToString(resultBig)
	require.Equal(t, expectedResult, am, "unexpected amount after add")
}

func TestSubAccountAmount(t *testing.T) {
	db := persistence.PostgresContext{
		Height: 0,
		DB:     *PostgresDB,
	}
	acc := NewTestAccount(t)

	err := db.SetAccountAmount(acc.Address, DefaultStake)
	require.NoError(t, err)

	amountToAddBig := big.NewInt(100)
	db.SubtractAccountAmount(acc.Address, types.BigIntToString(amountToAddBig))
	require.NoError(t, err)

	am, err := db.GetAccountAmount(acc.Address)
	require.NoError(t, err)

	resultBig := (&big.Int{}).Sub(DefaultStakeBig, amountToAddBig)
	expectedResult := types.BigIntToString(resultBig)
	require.Equal(t, expectedResult, am, "unexpected amount after sub")
}

// --- Pool Tests ---

func TestSetPoolAmount(t *testing.T) {
	db := persistence.PostgresContext{
		Height: 0,
		DB:     *PostgresDB,
	}
	pool := NewTestPool(t)

	err := db.SetPoolAmount(pool.Name, DefaultStake)
	require.NoError(t, err)

	am, err := db.GetPoolAmount(pool.Name)
	require.NoError(t, err)
	require.Equal(t, DefaultStake, am, "unexpected amount")

	err = db.SetPoolAmount(pool.Name, StakeToUpdate)
	require.NoError(t, err)

	am, err = db.GetPoolAmount(pool.Name)
	require.NoError(t, err)
	require.Equal(t, StakeToUpdate, am, "unexpected amount after second set")
}

func TestAddPoolAmount(t *testing.T) {
	db := persistence.PostgresContext{
		Height: 0,
		DB:     *PostgresDB,
	}
	pool := NewTestPool(t)

	err := db.SetPoolAmount(pool.Name, DefaultStake)
	require.NoError(t, err)

	amountToAddBig := big.NewInt(100)
	err = db.AddPoolAmount(pool.Name, types.BigIntToString(amountToAddBig))
	require.NoError(t, err)

	am, err := db.GetPoolAmount(pool.Name)
	require.NoError(t, err)

	resultBig := (&big.Int{}).Add(DefaultStakeBig, amountToAddBig)
	expectedResult := types.BigIntToString(resultBig)
	require.Equal(t, expectedResult, am, "unexpected amount after add")
}

func TestSubPoolAmount(t *testing.T) {
	db := persistence.PostgresContext{
		Height: 0,
		DB:     *PostgresDB,
	}
	pool := NewTestPool(t)

	err := db.SetPoolAmount(pool.Name, DefaultStake)
	require.NoError(t, err)

	amountToSubBig := big.NewInt(100)
	err = db.SubtractPoolAmount(pool.Name, types.BigIntToString(amountToSubBig))
	require.NoError(t, err)

	am, err := db.GetPoolAmount(pool.Name)
	require.NoError(t, err)

	resultBig := (&big.Int{}).Sub(DefaultStakeBig, amountToSubBig)
	expectedResult := types.BigIntToString(resultBig)
	require.Equal(t, expectedResult, am, "unexpected amount after sub")
}

// --- Helpers ---

func NewTestAccount(t *testing.T) typesGenesis.Account {
	addr, err := crypto.GenerateAddress()
	require.NoError(t, err)

	return typesGenesis.Account{
		Address: addr,
		Amount:  DefaultAccountAmount,
	}
}

func NewTestPool(t *testing.T) typesGenesis.Pool {
	_, err := crypto.GenerateAddress()
	require.NoError(t, err)

	return typesGenesis.Pool{
		Name: DefaultPoolName,
		Account: &typesGenesis.Account{
			Amount: DefaultAccountAmount,
		},
	}
}
