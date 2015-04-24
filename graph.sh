#!/bin/bash

OUT_DIR=$1

./plot_meanNearby.py $OUT_DIR/NearbyA.json $OUT_DIR/NearbyB.json &
./plot_meanGenetics.py $OUT_DIR/GeneticsA.json $OUT_DIR/GeneticsB.json &
./plot_stateSpace.py $OUT_DIR/Positions.json &
# ./plot_dead.py $OUT_DIR/DeadA.json $OUT_DIR/DeadB.json 2 &
./plot_Fitness.py $OUT_DIR/FitnessA.json $OUT_DIR/FitnessB.json 100 &
./draw.py $OUT_DIR/Positions.json