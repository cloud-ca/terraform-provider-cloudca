library(
  identifier: 'utils@v2.1.0',
  retriever: modernSCM([
    $class: 'GitSCMSource',
    remote: 'git@github.com:cloudops/cloudmc-jenkins-shared.git',
    credentialsId: 'gh-jenkins'
  ])
)

def cloudcaProviderRepo = 'terraform-provider-cloudca'
def targetBranch = 'master'
def releaseTypeName

properties([
    parameters([
        choice(name: 'BUMP', description: 'Which part of the version {major}.{minor}.{patch} to increase?',
            choices: [
                'major',
                'minor',
                'patch'
            ].join('\n'))
    ])
])


pipeline {
    agent {
        label 'cmc'
    }

    stages {
        stage('Setup'){
            steps {
                script {
                    releaseTypeName = params.BUMP    
                    checkout cloudcaProviderRepo targetBranch
                    sh 'git config user.name "jenkins"'
                    sh 'git config user.email "jenkins@cloudops.com"'
                }
            }
        }

        stage('Release') {
            steps {
                script {
                    sh "make ${releaseTypeName} push=true"
                }
            }
        }
    }
}

def checkout(repo, branch) {
    git url: "git@github.com:cloudca/$repo"+".git", branch: branch, credentialsId: 'gh-jenkins'
}