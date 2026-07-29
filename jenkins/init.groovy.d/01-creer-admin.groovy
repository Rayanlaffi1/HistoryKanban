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

if (realm.getUser(utilisateurAdmin) == null) {
  realm.createAccount(utilisateurAdmin, motDePasseAdmin)
  println "Administrateur Jenkins cree : ${utilisateurAdmin}"
} else {
  println "Administrateur Jenkins deja present : ${utilisateurAdmin}"
}

FullControlOnceLoggedInAuthorizationStrategy strategie = new FullControlOnceLoggedInAuthorizationStrategy()
strategie.setAllowAnonymousRead(false)
instance.setAuthorizationStrategy(strategie)
instance.save()
