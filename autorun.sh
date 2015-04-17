#!/bin/bash

go build

for i in {1..10}; do
	./iss
	mkdir ~/iss-data/a$i
	mv output/* ~/iss-data/a$i
	cp parameters.go ~/iss-data/a$i
	echo "done $i"
done