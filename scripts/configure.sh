#!/bin/bash

usage() {
    echo "usage: $0 [-h -s]"
    echo "  use -s to use build static libraries instead of dynamic one"
}

SETUP_MODS_OPTS=()

while getopts "hs" arg
do
    case "${arg}" in
        h)
            usage
            exit 0
            ;;
        s)
            SETUP_MODS_OPTS+=(--static-libraries)
            ;;
        *)
            usage
            exit 1
            ;;
    esac
done

echo -e "\e[0;1mBuilding tools...\e[0m"
make -C tools --no-print-directory

echo -e "\n\e[0;1mSetup modules...\e[0m"
./tools/bin/setup-modules --allow-overwrite ${SETUP_MODS_OPTS[@]}
echo -e "\n\e[0;1mSetup 3p modules...\e[0m"
./tools/bin/setup-3p-modules --allow-overwrite
