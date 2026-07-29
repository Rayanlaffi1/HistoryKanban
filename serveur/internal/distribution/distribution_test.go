package distribution

import (
	"reflect"
	"testing"

	"historykanban/serveur/internal/evenements"
)

func TestDestinatairesNotificationExclutActeurInteractif(t *testing.T) {
	e := evenements.Evenement{Acteur: "u1", Destinataires: []string{"u1", "u2", "u2"}}
	got := destinatairesNotification(e)
	want := []string{"u2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("destinataires = %#v, veut %#v", got, want)
	}
}

func TestDestinatairesNotificationInclutActeurAgent(t *testing.T) {
	e := evenements.Evenement{Acteur: "u1", Agent: true, Destinataires: []string{"u1", "u1"}}
	got := destinatairesNotification(e)
	want := []string{"u1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("destinataires = %#v, veut %#v", got, want)
	}
}
