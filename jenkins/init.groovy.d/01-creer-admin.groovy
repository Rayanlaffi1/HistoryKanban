import hudson.security.FullControlOnceLoggedInAuthorizationStrategy
import hudson.security.HudsonPrivateSecurityRealm
import jenkins.model.Jenkins

String utilisateurAdmin = System.getenv('JENKINSADMIN') ?: 'admin'
String motDePasseAdmin = System.getenv('JENKINSADMINMDP') ?: ''

if (!motDePasseAdmin.trim()) {
  println 'Creation de l administrateur Jenkins ignoree : JENKINSADMINMDP est vide.'
  return
}

Jenkins instance = Jenkins.get()
HudsonPrivateSecurityRealm realm

if (instance.getSecurityRealm() instanceof HudsonPrivateSecurityRealm) {
  realm = (HudsonPrivateSecurityRealm) instance.getSecurityRealm()
} else {
  realm = new HudsonPrivateSecurityRealm(false)
  instance.setSecurityRealm(realm)
}

def utilisateur = realm.getUser(utilisateurAdmin)
if (utilisateur == null) {
  utilisateur = realm.createAccount(utilisateurAdmin, motDePasseAdmin)
  println "Administrateur Jenkins cree : ${utilisateurAdmin}"
} else {
  utilisateur.addProperty(HudsonPrivateSecurityRealm.Details.fromPlainPassword(motDePasseAdmin))
  utilisateur.save()
  println "Mot de passe de l administrateur Jenkins synchronise : ${utilisateurAdmin}"
}

FullControlOnceLoggedInAuthorizationStrategy strategie = new FullControlOnceLoggedInAuthorizationStrategy()
strategie.setAllowAnonymousRead(false)
instance.setAuthorizationStrategy(strategie)
instance.save()
