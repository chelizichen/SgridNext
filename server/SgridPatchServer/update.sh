#!/usr/bin/env bash
# Written in [Amber](https://amber-lang.com/)
# version: 0.3.5-alpha
# date: 2025-09-01 15:34:05


exit__80_v0() {
    local code=$1
    exit "${code}";
    __AS=$?
}
args=("$0" "$@")
    node_dir="${args[1]}"
    echo "node_dir: ${node_dir}"
    echo "update sgridnode"
    echo "rm ${node_dir}/sgridnode"
     rm ${node_dir}/sgridnode ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "rm ${node_dir}/sgridnode failed"
        exit__80_v0 1;
        __AF_exit80_v0__11_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__11_9" > /dev/null 2>&1
fi
    echo "cp ./sgridnode ${node_dir}"
     cp ./sgridnode ${node_dir} ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cp ./sgridnode ${node_dir} failed"
        exit__80_v0 1;
        __AF_exit80_v0__16_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__16_9" > /dev/null 2>&1
fi
    echo "chmod +x ${node_dir}/sgridnode"
     chmod +x ${node_dir}/sgridnode ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "chmod +x ${node_dir}/sgridnode failed"
        exit__80_v0 1;
        __AF_exit80_v0__21_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__21_9" > /dev/null 2>&1
fi
    echo "systemctl restart sgridnode"
     systemctl restart sgridnode ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "systemctl restart sgridnode failed"
        exit__80_v0 1;
        __AF_exit80_v0__26_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__26_9" > /dev/null 2>&1
fi
    echo "update success"
