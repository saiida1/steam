#!/bin/bash
gnome-terminal -x sh -c "./gradlew jettyRunWar" ; sleep 10 ; curl -X POST --form pojo=@examples/pojo-server/gbm_3f258f27_f0ad_4520_b6a5_3d2bb4a9b0ff.java --form jar=@examples/pojo-server/h2o-genmodel.jar localhost:55000/makewar > example.war 
