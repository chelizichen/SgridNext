#!/usr/bin/env bash
# Written in [Amber](https://amber-lang.com/)
# version: 0.3.5-alpha
# date: 2025-09-01 15:03:30


exit__80_v0() {
    local code=$1
    exit "${code}";
    __AS=$?
}
args=("$0" "$@")
     rm ./sgridnode ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "rm ./sgridnode failed, ignore"
        # exit(1)
fi
     rm ./sgridnode.tar.gz ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "rm ./sgridnode.tar.gz failed, ignore"
        # exit(1)
fi
     GOOS=linux GOARCH=amd64 go build -o update;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "go build failed"
        exit__80_v0 1;
        __AF_exit80_v0__16_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__16_9" > /dev/null 2>&1
fi
     cd ../../;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd ../../ failed"
        exit__80_v0 1;
        __AF_exit80_v0__21_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__21_9" > /dev/null 2>&1
fi
     sh deploy.sh ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "sh deploy.sh failed"
        exit__80_v0 1;
        __AF_exit80_v0__26_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__26_9" > /dev/null 2>&1
fi
     cp ./archive/sgridnode/sgridnode ./server/SgridPatchServer ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cp ./archive/sgridnode/sgridnode ./server/SgridPatchServer failed"
        exit__80_v0 1;
        __AF_exit80_v0__31_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__31_9" > /dev/null 2>&1
fi
     cd server/SgridPatchServer ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd server/SgridPatchServer failed"
        exit__80_v0 1;
        __AF_exit80_v0__36_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__36_9" > /dev/null 2>&1
fi
     tar -zcvf sgridnode.tar.gz ./* ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "tar -zcvf sgridnode.tar.gz ./* failed"
        exit__80_v0 1;
        __AF_exit80_v0__41_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__41_9" > /dev/null 2>&1
fi
    echo "deploy success"
    exit__80_v0 0;
    __AF_exit80_v0__44_5="$__AF_exit80_v0";
    echo "$__AF_exit80_v0__44_5" > /dev/null 2>&1
