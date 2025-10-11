#!/usr/bin/env bash
# Written in [Amber](https://amber-lang.com/)
# version: 0.3.5-alpha
# date: 2025-10-11 10:25:33


exit__80_v0() {
    local code=$1
    exit "${code}";
    __AS=$?
}
__0_SGRID_NODE="sgridnode"
__1_SGRID_NEXT="sgridnext"
__2_PKG="SgridCloud"
clean__94_v0() {
    rm -r dist;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "remove web dist failed, ignore"
        # exit(1)
fi
    rm -r archive;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "remove archive failed, ignore"
        # exit(1)
fi
}
build_sgridnext_web__95_v0() {
     cd web ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd web failed"
        exit__80_v0 1;
        __AF_exit80_v0__21_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__21_9" > /dev/null 2>&1
fi
     npm run build ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "npm run build failed"
        exit__80_v0 1;
        __AF_exit80_v0__25_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__25_9" > /dev/null 2>&1
fi
     cd .. ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd .. failed"
        exit__80_v0 1;
        __AF_exit80_v0__29_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__29_9" > /dev/null 2>&1
fi
     cp -r web/dist dist ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cp -r web/dist dist failed"
        exit__80_v0 1;
        __AF_exit80_v0__33_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__33_9" > /dev/null 2>&1
fi
}
build_sgridnext_backend__96_v0() {
    GOOS=linux GOARCH=amd64 go build -o "${__1_SGRID_NEXT}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "build sgridnext backend failed"
        exit__80_v0 1;
        __AF_exit80_v0__41_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__41_9" > /dev/null 2>&1
fi
}
build_sgrid_node__97_v0() {
     cd server/SgridNodeServer ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd server/SgridNodeServer failed"
        exit__80_v0 1;
        __AF_exit80_v0__48_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__48_9" > /dev/null 2>&1
fi
     GOOS=linux GOARCH=amd64 go build -o "${__0_SGRID_NODE}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "build sgrid node backend failed"
        exit__80_v0 1;
        __AF_exit80_v0__52_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__52_9" > /dev/null 2>&1
fi
     mv ./"${__0_SGRID_NODE}" ../../ ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "mv sgrid node backend failed"
        exit__80_v0 1;
        __AF_exit80_v0__57_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__57_9" > /dev/null 2>&1
fi
     cd ../../ ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd../../ failed"
        exit__80_v0 1;
        __AF_exit80_v0__61_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__61_9" > /dev/null 2>&1
fi
}
make_archive__98_v0() {
     mkdir archive ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "mkdir archive failed"
        exit__80_v0 1;
        __AF_exit80_v0__68_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__68_9" > /dev/null 2>&1
fi
     mkdir archive/"${__1_SGRID_NEXT}"  ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "mkdir archive/"
                    echo "${__1_SGRID_NEXT}" > /dev/null 2>&1
        echo " failed" > /dev/null 2>&1
        exit__80_v0 1;
        __AF_exit80_v0__72_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__72_9" > /dev/null 2>&1
fi
     mkdir archive/"${__0_SGRID_NODE}"  ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "mkdir archive/"
                    echo "${__0_SGRID_NODE}" > /dev/null 2>&1
        echo " failed" > /dev/null 2>&1
        exit__80_v0 1;
        __AF_exit80_v0__76_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__76_9" > /dev/null 2>&1
fi
     mv ./"${__0_SGRID_NODE}" ./archive/"${__0_SGRID_NODE}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "mv sgrid node backend failed"
        exit__80_v0 1;
        __AF_exit80_v0__81_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__81_9" > /dev/null 2>&1
fi
     cp ./server/SgridNodeServer/config.template.json ./archive/"${__0_SGRID_NODE}"/config.json ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cp sgrid next backend failed"
        exit__80_v0 1;
        __AF_exit80_v0__85_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__85_9" > /dev/null 2>&1
fi
     mv ./"${__1_SGRID_NEXT}" ./archive/"${__1_SGRID_NEXT}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "mv sgrid next backend failed"
        exit__80_v0 1;
        __AF_exit80_v0__90_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__90_9" > /dev/null 2>&1
fi
     cp ./config.json ./archive/"${__1_SGRID_NEXT}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cp sgrid next backend failed"
        exit__80_v0 1;
        __AF_exit80_v0__94_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__94_9" > /dev/null 2>&1
fi
     mv ./dist ./archive/"${__1_SGRID_NEXT}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "mv web dist failed"
        exit__80_v0 1;
        __AF_exit80_v0__99_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__99_9" > /dev/null 2>&1
fi
     cp sgridnext.service ./archive ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cp sgridnext.service failed"
        exit__80_v0 1;
        __AF_exit80_v0__104_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__104_9" > /dev/null 2>&1
fi
     cp sgridnode.service ./archive ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cp sgridnode.service failed"
        exit__80_v0 1;
        __AF_exit80_v0__109_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__109_9" > /dev/null 2>&1
fi
}
tar_archive__99_v0() {
     tar -zcvf "${__2_PKG}".tar.gz archive ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "tar -zcvf "
                    echo "${__2_PKG}" > /dev/null 2>&1
        echo ".tar.gz archive failed" > /dev/null 2>&1
        exit__80_v0 1;
        __AF_exit80_v0__116_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__116_9" > /dev/null 2>&1
fi
     cd archive/"${__0_SGRID_NODE}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd archive/"
                    echo "${__0_SGRID_NODE}" > /dev/null 2>&1
        echo " failed" > /dev/null 2>&1
        exit__80_v0 1;
        __AF_exit80_v0__120_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__120_9" > /dev/null 2>&1
fi
    tar -zcvf "${__0_SGRID_NODE}".tar.gz ./*;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "tar -zcvf "
                    echo "${__0_SGRID_NODE}" > /dev/null 2>&1
        echo ".tar.gz archive failed" > /dev/null 2>&1
        exit__80_v0 1;
        __AF_exit80_v0__124_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__124_9" > /dev/null 2>&1
fi
     cd ../../ ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd ../../ failed"
        exit__80_v0 1;
        __AF_exit80_v0__129_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__129_9" > /dev/null 2>&1
fi
     cd archive/"${__1_SGRID_NEXT}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "cd archive/"
                    echo "${__1_SGRID_NEXT}" > /dev/null 2>&1
        echo " failed" > /dev/null 2>&1
        exit__80_v0 1;
        __AF_exit80_v0__134_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__134_9" > /dev/null 2>&1
fi
    tar -zcvf "${__1_SGRID_NEXT}".tar.gz ./* ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "tar -zcvf "
                    echo "${__1_SGRID_NEXT}" > /dev/null 2>&1
        echo ".tar.gz archive failed" > /dev/null 2>&1
        exit__80_v0 1;
        __AF_exit80_v0__139_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__139_9" > /dev/null 2>&1
fi
}
args=("$0" "$@")
    clean__94_v0 ;
    __AF_clean94_v0__144_5="$__AF_clean94_v0";
    echo "$__AF_clean94_v0__144_5" > /dev/null 2>&1
    build_sgridnext_web__95_v0 ;
    __AF_build_sgridnext_web95_v0__145_5="$__AF_build_sgridnext_web95_v0";
    echo "$__AF_build_sgridnext_web95_v0__145_5" > /dev/null 2>&1
    build_sgridnext_backend__96_v0 ;
    __AF_build_sgridnext_backend96_v0__146_5="$__AF_build_sgridnext_backend96_v0";
    echo "$__AF_build_sgridnext_backend96_v0__146_5" > /dev/null 2>&1
    build_sgrid_node__97_v0 ;
    __AF_build_sgrid_node97_v0__147_5="$__AF_build_sgrid_node97_v0";
    echo "$__AF_build_sgrid_node97_v0__147_5" > /dev/null 2>&1
    make_archive__98_v0 ;
    __AF_make_archive98_v0__148_5="$__AF_make_archive98_v0";
    echo "$__AF_make_archive98_v0__148_5" > /dev/null 2>&1
    tar_archive__99_v0 ;
    __AF_tar_archive99_v0__149_5="$__AF_tar_archive99_v0";
    echo "$__AF_tar_archive99_v0__149_5" > /dev/null 2>&1
