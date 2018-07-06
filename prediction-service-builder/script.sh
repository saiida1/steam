#!/bin/bash
#ps -aux | grep gradlew
export JAVA_TOOL_OPTIONS=-Dfile.encoding=UTF8
export _JAVA_OPTIONS=-Xmx12096m
(./gradlew jettyRunWar & sleep 10) ; curl -X POST --form pojo=@/home/saida/pojos/rf.java --form jar=@/home/saida/pojos/h2o-genmodel.jar localhost:55000/makewar > example.war ; sleep 120 
exit 0
