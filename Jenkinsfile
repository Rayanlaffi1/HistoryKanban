pipeline {
  agent any

  options {
    timestamps()
    disableConcurrentBuilds()
  }

  stages {
    stage("Preparation") {
      steps {
        checkout scm
      }
    }

    stage("Verification serveur") {
      steps {
        dir("serveur") {
          sh "docker build --target verification -t historykanban-serveur-verification ."
        }
      }
    }

    stage("Construction serveur") {
      steps {
        dir("serveur") {
          sh "docker build -t historykanban-serveur:${env.BUILD_NUMBER} -t historykanban-serveur:latest ."
        }
      }
    }

    stage("Construction interface") {
      steps {
        dir("interface") {
          sh "docker build -t historykanban-interface:${env.BUILD_NUMBER} -t historykanban-interface:latest ."
        }
      }
    }

    stage("Deploiement") {
      when {
        anyOf {
          branch "main"
          expression { env.BRANCH_NAME == null }
        }
      }
      steps {
        sh "docker compose up -d --build"
      }
    }
  }

  post {
    success {
      echo "Deploiement HistoryKanban reussi"
    }
    failure {
      echo "Echec du pipeline HistoryKanban"
    }
  }
}
