#!/bin/bash

./plot_meanNearby.py output/NearbyA.json output/NearbyB.json 4000 2 &
./plot_meanGenetics.py output/GeneticsA.json output/GeneticsB.json -alt &
./plot_dead.py output/DeadA.json output/DeadB.json 2 &
./draw.py output/Positions.json