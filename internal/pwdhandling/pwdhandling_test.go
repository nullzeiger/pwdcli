package pwdhandling

import (
	"os"
	"testing"

	"github.com/nullzeiger/pwdcli/internal/filehandling"
)

func withTempHome(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()
	if err := os.Setenv("HOME", tmp); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	return tmp
}

func writeAccounts(t *testing.T, acc []filehandling.Account) {
	t.Helper()
	if err := filehandling.WriteJSON(acc); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}
}

func TestAll(t *testing.T) {
	tmp := withTempHome(t)

	accounts := []filehandling.Account{
		{Website: "example.com", Username: "john", Email: "j@example.com", Pwd: "123"},
		{Website: "test.com", Username: "alice", Email: "a@test.com", Pwd: "456"},
	}
	writeAccounts(t, accounts)

	list, err := All()
	if err != nil {
		t.Fatalf("All() failed: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}

	expected := "[0] Website: example.com Username: john Email: j@example.com Password: 123"
	if list[0] != expected {
		t.Fatalf("unexpected first entry:\n%s", list[0])
	}

	_ = tmp
}

func TestCreate(t *testing.T) {
	tmp := withTempHome(t)

	if err := filehandling.Create(); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	acc := filehandling.Account{
		Website:  "site.com",
		Username: "user",
		Email:    "u@e.com",
		Pwd:      "pwd",
	}

	if err := Create(acc); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	read, err := filehandling.ReadJSON()
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if len(read) != 1 {
		t.Fatalf("expected 1 account, got %d", len(read))
	}

	if read[0].Website != "site.com" {
		t.Fatalf("unexpected account: %+v", read[0])
	}

	_ = tmp
}

func TestDelete(t *testing.T) {
	tmp := withTempHome(t)

	accounts := []filehandling.Account{
		{Website: "1.com", Username: "a", Email: "a@a", Pwd: "1"},
		{Website: "2.com", Username: "b", Email: "b@b", Pwd: "2"},
		{Website: "3.com", Username: "c", Email: "c@c", Pwd: "3"},
	}
	writeAccounts(t, accounts)

	ok, err := Delete(1)
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if !ok {
		t.Fatalf("Delete returned false when expected true")
	}

	read, err := filehandling.ReadJSON()
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	if len(read) != 2 {
		t.Fatalf("expected length 2 after deletion, got %d", len(read))
	}

	if read[1].Website != "3.com" {
		t.Fatalf("unexpected remaining entries: %+v", read)
	}

	_, err = Delete(10)
	if err == nil {
		t.Fatalf("expected error for out-of-range index")
	}

	_ = tmp
}

func TestSearch(t *testing.T) {
	tmp := withTempHome(t)

	accounts := []filehandling.Account{
		{Website: "google.com", Username: "john", Email: "john@google.com", Pwd: "xyz"},
		{Website: "github.com", Username: "alice", Email: "alice@gh.com", Pwd: "pwd"},
		{Website: "example.org", Username: "bob", Email: "b@ex.org", Pwd: "123"},
	}
	writeAccounts(t, accounts)

	results, err := Search("git")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 search result, got %d", len(results))
	}

	if results[0].Account.Website != "github.com" {
		t.Fatalf("unexpected search result: %+v", results[0].Account)
	}

	results, err = Search("JOHN")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 || results[0].Account.Username != "john" {
		t.Fatalf("case-insensitive search failed")
	}

	_ = tmp
}

func TestAllEmptyFile(t *testing.T) {
	tmp := withTempHome(t)

	writeAccounts(t, []filehandling.Account{})

	entries, err := All()
	if err != nil {
		t.Fatalf("All() error: %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected empty result, got %v", entries)
	}

	_ = tmp
}
