package bank

import "testing"

func TestRandomSessionBankLifecycle(t *testing.T) {
	bank := NewRandomSessionBank()
	session := &RandomSession{ID: 8, Questions: []int{1, 2}}
	if err := bank.Create(session); err != nil {
		t.Fatalf("create: %v", err)
	}
	found, err := bank.FindByID(8)
	if err != nil || found != session {
		t.Fatalf("unexpected found session: %+v, %v", found, err)
	}
	session.CurrentIndex = 1
	if err := bank.Save(session); err != nil {
		t.Fatalf("save: %v", err)
	}
	found, _ = bank.FindByID(8)
	if found.CurrentIndex != 1 {
		t.Fatalf("session update not saved: %+v", found)
	}
	if err := bank.Delete(8); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := bank.FindByID(8); err == nil {
		t.Fatal("expected not found after delete")
	}
}
