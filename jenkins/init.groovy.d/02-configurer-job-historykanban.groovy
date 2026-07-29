import hudson.plugins.git.BranchSpec
import hudson.plugins.git.GitSCM
import jenkins.model.Jenkins
import org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition
import org.jenkinsci.plugins.workflow.job.WorkflowJob

Jenkins instance = Jenkins.get()
String nomJob = System.getenv('JENKINSJOB') ?: 'historykanban'
String depotGit = System.getenv('JENKINSGITURL') ?: 'https://github.com/Rayanlaffi1/HistoryKanban.git'
String branche = System.getenv('JENKINSBRANCHE') ?: 'main'

if (nomJob.contains('/')) {
  println "Creation automatique du job Jenkins ignoree : JENKINSJOB ne doit pas contenir de dossier pour l'initialisation simple (${nomJob})."
  return
}

WorkflowJob job = instance.getItem(nomJob) as WorkflowJob
if (job == null) {
  job = instance.createProject(WorkflowJob, nomJob)
  println "Job Jenkins cree : ${nomJob}"
} else {
  println "Job Jenkins deja present : ${nomJob}"
}

GitSCM scm = new GitSCM(depotGit)
scm.branches = [new BranchSpec("*/${branche}")]

CpsScmFlowDefinition definition = new CpsScmFlowDefinition(scm, 'Jenkinsfile')
definition.setLightweight(true)
job.setDefinition(definition)
job.save()

println "Job Jenkins ${nomJob} configure sur ${depotGit} (${branche})."
