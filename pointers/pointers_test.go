package pointers

import "testing"

func TestPointer(t *testing.T) {
	assertBalance := func(got Bitcoin, want Bitcoin) {
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		got := wallet.Balance()
		want := Bitcoin(10)

		assertBalance(got, want)
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		wallet.Withdraw(Bitcoin(12))
		got := wallet.Balance()
		want := Bitcoin(8)

		assertBalance(got, want)
	})
}