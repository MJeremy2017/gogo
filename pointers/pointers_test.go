package pointers

import "testing"

func TestPointer(t *testing.T) {
	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		got := wallet.Balance()
		want := Bitcoin(10)

		assertBalance(t, got, want)
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(12))
		got := wallet.Balance()
		want := Bitcoin(8)

		assertBalance(t, got, want)
		if err != nil {

		}
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(10)}
		err := wallet.Withdraw(20)
		got := wallet.Balance()
		want := Bitcoin(10)

		assertBalance(t, got, want)

		assertError(t, err, ErrInsufficientFunds)

	})
}

func assertBalance(t testing.TB, got Bitcoin, want Bitcoin) {
		t.Helper()
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	}

func assertError(t testing.TB, err error, want error) {
	t.Helper()
	if err == nil {
		// fatal: stop if reached here
		t.Fatal("Intend to receive error but not!")
	}

	if err != want {
		t.Errorf("err wanted %q got %q", want, err)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Error("Got an error but didn't want one!")
	}
}




