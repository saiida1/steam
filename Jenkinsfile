#!/usr/bin/env groovy
def function(String dockerImageName){
    stage('checkout') {
        checkout scm
    }
    stage('clean') {
    dir('prediction-service-builder')  {
    sh './gradlew clean build'
    }
                }
  stage('Build model') {
      dir('prediction-service-builder')  {
        sh './gradlew build'
        sh 'export JAVA_TOOL_OPTIONS=-Dfile.encoding=UTF8'
        sh 'export _JAVA_OPTIONS=-Xmx12096m'
       // sh 'gnome-terminal -x sh -c "./gradlew jettyRunWar" ; sleep 10 ; curl -X POST --form pojo=@examples/pojo-server/gbm_3f258f27_f0ad_4520_b6a5_3d2bb4a9b0ff.java --form jar=@examples/pojo-server/h2o-genmodel.jar localhost:55000/makewar > example.war'
    sh'(./gradlew jettyRunWar & sleep 10) ; curl -X POST --form pojo=@/home/saida/pojos/rf.java --form jar=@/home/saida/pojos/h2o-genmodel.jar localhost:55000/makewar > example.war ; sleep 60 '
      }
  }
 
  stage('packaging') {
    dir('prediction-service-builder')  {
    sh 'mv target/**/* target'   
    dir('target')  {
    sh 'jar -cvf prediction.war *'
    }
    }
}
}
    node{

      dir('prediction-service') {
        stage('predictionservice') {
          echo "** prediction service micro-service ***"
        }

        function('predictionservice')
      }
    }
