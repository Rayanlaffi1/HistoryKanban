package configuration

import (
	"os"
	"strings"
)

type Configuration struct {
	Port                string
	BDURL               string
	KeycloakURL         string
	KeycloakURLPublique string
	KeycloakRealm       string
	MinioHote           string
	MinioURLPublique    string
	MinioCle            string
	MinioSecret         string
	MinioSeau           string
	RabbitURL           string
	SMTPHote            string
	SMTPPort            string
	SMTPExpediteur      string
	Origines            []string
	JenkinsURL          string
	JenkinsJob          string
	JenkinsUtilisateur  string
	JenkinsJeton        string
	GithubDepot         string
	VersionApplication  string
}

func lire(cle, defaut string) string {
	valeur := os.Getenv(cle)
	if valeur == "" {
		return defaut
	}
	return valeur
}

func Charger() Configuration {
	return Configuration{
		Port:                lire("PORTSERVEUR", "8080"),
		BDURL:               lire("BDURL", "postgres://historykanban:historykanban@localhost:5432/historykanban?sslmode=disable"),
		KeycloakURL:         lire("KEYCLOAKURL", "http://localhost:8081"),
		KeycloakURLPublique: lire("KEYCLOAKURLPUBLIQUE", "https://auth.historykanban.localhost"),
		KeycloakRealm:       lire("KEYCLOAKREALM", "historykanban"),
		MinioHote:           lire("MINIOHOTE", "localhost:9000"),
		MinioURLPublique:    lire("MINIOURLPUBLIQUE", "https://images.historykanban.localhost"),
		MinioCle:            lire("MINIOCLE", "historykanban"),
		MinioSecret:         lire("MINIOSECRET", "historykanban"),
		MinioSeau:           lire("MINIOSEAU", "images"),
		RabbitURL:           lire("RABBITURL", "amqp://historykanban:historykanban@localhost:5672/"),
		SMTPHote:            lire("SMTPHOTE", "localhost"),
		SMTPPort:            lire("SMTPPORT", "1025"),
		SMTPExpediteur:      lire("SMTPEXPEDITEUR", "notifications@historykanban.fr"),
		Origines:            strings.Split(lire("ORIGINES", "http://localhost:5173"), ","),
		JenkinsURL:          lire("JENKINSURL", ""),
		JenkinsJob:          lire("JENKINSJOB", ""),
		JenkinsUtilisateur:  lire("JENKINSUTILISATEUR", ""),
		JenkinsJeton:        lire("JENKINSJETON", ""),
		GithubDepot:         lire("GITHUBDEPOT", ""),
		VersionApplication:  lire("VERSIONAPPLICATION", "dev"),
	}
}

func (c Configuration) URLJWKS() string {
	return c.KeycloakURL + "/realms/" + c.KeycloakRealm + "/protocol/openid-connect/certs"
}

// Emetteurs accepte une liste d'URL publiques separees par des virgules,
// pour valider les jetons emis via historykanban.localhost comme via un
// domaine reseau du type historykanban.<IP>.sslip.io.
func (c Configuration) Emetteurs() []string {
	var emetteurs []string
	for _, base := range strings.Split(c.KeycloakURLPublique, ",") {
		base = strings.TrimRight(strings.TrimSpace(base), "/")
		if base != "" {
			emetteurs = append(emetteurs, base+"/realms/"+c.KeycloakRealm)
		}
	}
	return emetteurs
}
