pipeline {
    agent any
    tools {
        go 'go-1.23.3'
    }
    triggers {
        githubPush()
    }
    stages {
        stage('Checkout Code') {
            steps {
                checkout scm
            }
        }
        stage('Run Go Tests') {
            steps {
                script {
                    def goModules = sh(script: 'find . -name "go.mod" -exec dirname {} \\;', returnStdout: true).trim().split("\n")
                    
                    goModules.each { module ->
                        echo "Running tests for module: ${module}"
                        sh "cd ${module} && go test -v ./..."
                    }
                }
            }
        }
    }
    post {
        always {
            cleanWs()
        }
        success {
            script {
                githubNotify(
                    status: 'SUCCESS',
                    description: 'Build completed successfully.',
                    credentialsId: 'github-testoken',
                    context: 'ci/build', 
                    targetUrl: "http://cams-jenkins.duckdns.org:50000/job/K8sDemoUnitTests/${env.BUILD_NUMBER}", 
                    repo: 'K8sDemo',
                    account: 'ACamaraLara',  
                    sha: env.GIT_COMMIT  // SHA del commit actual
                )
            }
        }
        failure {
            script {
                githubNotify(
                    status: 'FAILURE',
                    description: 'Build failed.',
                    credentialsId: 'github-testoken',
                    context: 'ci/build',
                    targetUrl: "http://cams-jenkins.duckdns.org:50000/job/K8sDemoUnitTests/${env.BUILD_NUMBER}",
                    repo: 'K8sDemo',
                    account: 'ACamaraLara',
                    sha: env.GIT_COMMIT
                )
            }
        }
    }
}

