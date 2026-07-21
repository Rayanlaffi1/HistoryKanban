package courriels

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Envoyeur struct {
	adresse    string
	expediteur string
}

func Nouveau(hote, port, expediteur string) *Envoyeur {
	return &Envoyeur{adresse: hote + ":" + port, expediteur: expediteur}
}

func (e *Envoyeur) Envoyer(destinataires []string, sujet, corps string) error {
	if len(destinataires) == 0 {
		return nil
	}
	message := strings.Join([]string{
		"From: HistoryKanban <" + e.expediteur + ">",
		"To: " + strings.Join(destinataires, ", "),
		"Subject: " + sujet,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		corps,
	}, "\r\n")
	return smtp.SendMail(e.adresse, nil, e.expediteur, destinataires, []byte(message))
}

func (e *Envoyeur) GabaritChangement(titre, action, acteur, projet string) string {
	return fmt.Sprintf(`<div style="font-family:sans-serif;background:#f5f5f5;padding:24px">
<div style="background:#ffffff;border:1px solid #d4d4d4;border-radius:12px;padding:24px;max-width:520px;margin:auto">
<h2 style="color:#171717;margin:0 0 16px">HistoryKanban</h2>
<p style="color:#525252">%s</p>
<p style="color:#171717;font-weight:bold">%s</p>
<p style="color:#525252">Par %s — projet %s</p>
</div></div>`, action, titre, acteur, projet)
}
