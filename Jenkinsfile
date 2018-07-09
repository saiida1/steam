pipeline {
    agent any
    triggers {
      cron('H 1 * * *')
    }
    stages {
        stage('checkout') {
            steps {
                checkout scm
            }
        }
        stage('clean') {
            steps {
   		        dir('prediction-service-builder') {
                sh './gradlew clean build'
              }
            }
        }

        stage('Build model') {
            steps {
   		        dir('prediction-service-builder') {
    		        sh './gradlew build'
    			      sh 'chmod +x script.sh'
    			      sh './script.sh'
              }
	          }
	      }
        stage('packaging') {
            steps {
   		        dir('prediction-service-builder') {
    		        sh 'mv target/**/* target'
                dir('target') {
      		        sh 'rm -r makeWar*'
      		        sh 'jar -cvf prediction.war *'
                }
                stash includes : 'target/*.war', name: 'prediction.war'
  	          }
	          }
        }
}
}
