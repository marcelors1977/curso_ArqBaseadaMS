#/bin/bash

awk -F':' '{ system("kafka-topics --create --topic="$1" --if-not-exists --bootstrap-server=kafka:29092" ) }' all_topics.txt