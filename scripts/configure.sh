#!/bin/bash

echo -e "\e[0;1mBuilding tools...\e[0m"
make -C tools --no-print-directory

echo -e "\n\e[0;1mSetup modules...\e[0m"
./tools/bin/setup-modules --allow-overwrite
echo -e "\n\e[0;1mSetup 3p modules...\e[0m"
./tools/bin/setup-3p-modules --allow-overwrite
