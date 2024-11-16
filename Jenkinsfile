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
            echo 'Tests passed successfully!'
        }
        failure {
            echo 'Tests failed!'
        }
    }
}
