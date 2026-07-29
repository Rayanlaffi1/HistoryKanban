package depots

import "testing"

func TestEmpreinteCleHashSHA256(t *testing.T) {
	emp := empreinteCle("hk.secret")
	if emp == "hk.secret" {
		t.Fatal("la cle ne doit pas rester en clair")
	}
	if len(emp) != 64 || !estEmpreinteCle(emp) {
		t.Fatalf("empreinte invalide: %q", emp)
	}
	if emp != empreinteCle("hk.secret") {
		t.Fatal("le hash doit etre stable pour authentifier la cle")
	}
}
