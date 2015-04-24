#!/bin/bash

go build

for i in {0..50}; do
	tmpDir=$(mktemp -d)
	baseTmpDir=$(basename $tmpDir)

	# copy binary to new folder
	cp ~/go/src/github.com/oskanberg/iss/iss $tmpDir

	# cd there, run script
	pushd $tmpDir
	mkdir output
	./iss

	# create results dir, copy
	mkdir ~/iss-data/assym-$baseTmpDir
	mv output/* ~/iss-data/assym-$baseTmpDir

	# cleanup 
	popd
	rm -rf $tmpDir

	echo "done $i"
done