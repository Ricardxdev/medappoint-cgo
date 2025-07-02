#!/bin/bash

cd csrc
gcc -I./ -o main patient.c patient_metrics.c
mv main ../
cd ../
./main
