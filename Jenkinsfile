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
    //    sh './script.sh'
       // sh 'gnome-terminal -x sh -c "./gradlew jettyRunWar" ; sleep 10 ; curl -X POST --form pojo=@examples/pojo-server/gbm_3f258f27_f0ad_4520_b6a5_3d2bb4a9b0ff.java --form jar=@examples/pojo-server/h2o-genmodel.jar localhost:55000/makewar > example.war'
    sh'  (./gradlew jettyRunWar & sleep 10) ; curl -X POST --form pojo=@examples/pojo-server/gbm_3f258f27_f0ad_4520_b6a5_3d2bb4a9b0ff.java --form jar=@examples/pojo-server/h2o-genmodel.jar localhost:55000/makewar > example.war'
      }
  }
 
  stage('packaging') {
    dir('target')  {
    sh 'jar -cvf prediction.war *'
    stash includes : 'prediction.war/*', name: 'prediction'
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
