#!/usr/bin/env bash
# Written in [Amber](https://amber-lang.com/)
# version: 0.3.5-alpha
# date: 2025-10-10 09:57:06
dir_exist__0_v0() {
    local path=$1
    [ -d "${path}" ];
    __AS=$?;
if [ $__AS != 0 ]; then
        __AF_dir_exist0_v0=0;
        return 0
fi
    __AF_dir_exist0_v0=1;
    return 0
}
file_exist__1_v0() {
    local path=$1
    [ -f "${path}" ];
    __AS=$?;
if [ $__AS != 0 ]; then
        __AF_file_exist1_v0=0;
        return 0
fi
    __AF_file_exist1_v0=1;
    return 0
}
file_write__3_v0() {
    local path=$1
    local content=$2
    __AMBER_VAL_0=$(echo "${content}" > "${path}");
    __AS=$?;
if [ $__AS != 0 ]; then
__AF_file_write3_v0=''
return $__AS
fi;
    __AF_file_write3_v0="${__AMBER_VAL_0}";
    return 0
}
create_dir__6_v0() {
    local path=$1
    dir_exist__0_v0 "${path}";
    __AF_dir_exist0_v0__48_12="$__AF_dir_exist0_v0";
    if [ $(echo  '!' "$__AF_dir_exist0_v0__48_12" | bc -l | sed '/\./ s/\.\{0,1\}0\{1,\}$//') != 0 ]; then
        mkdir -p "${path}";
        __AS=$?
fi
}

exit__80_v0() {
    local code=$1
    exit "${code}";
    __AS=$?
}
create_dir_if_not_exists__99_v0() {
    local path=$1
    dir_exist__0_v0 "${path}";
    __AF_dir_exist0_v0__5_13="$__AF_dir_exist0_v0";
    if [ $(echo  '!' "$__AF_dir_exist0_v0__5_13" | bc -l | sed '/\./ s/\.\{0,1\}0\{1,\}$//') != 0 ]; then
        create_dir__6_v0 "${path}";
        __AF_create_dir6_v0__6_9="$__AF_create_dir6_v0";
        echo "$__AF_create_dir6_v0__6_9" > /dev/null 2>&1
fi
}
create_service_file__100_v0() {
    local install_dir=$1
    local service_path="/usr/lib/systemd/system/sgridnode.service"
     touch "${service_path}" ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "create sgridnode.service faild"
        exit__80_v0 1;
        __AF_exit80_v0__14_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__14_9" > /dev/null 2>&1
fi
    local content="
    [Unit]
    Description = sgrid next,A cloud platform for grid computing
    [Service]
    Type = simple
    ExecStart = ${install_dir}/sgridnode
    WorkingDirectory = ${install_dir}
    Environment=PATH=/usr/bin:/usr/local/bin
    Restart = no
    "
    file_write__3_v0 "${service_path}" "${content}";
    __AS=$?;
    __AF_file_write3_v0__26_12="${__AF_file_write3_v0}";
    echo "${__AF_file_write3_v0__26_12}" > /dev/null 2>&1
}
# sh install.sh 10.124.13.111 /home/sgridnode/ sgridnode.tar.gz
args=("$0" "$@")
    host="${args[1]}"
    install_dir="${args[2]}"
    package_path="${args[3]}"
    echo "install host -> ${host}"
    echo "install install_dir -> ${install_dir}"
    # 检查安装包是否存在
    file_exist__1_v0 "${package_path}";
    __AF_file_exist1_v0__39_13="$__AF_file_exist1_v0";
    if [ $(echo  '!' "$__AF_file_exist1_v0__39_13" | bc -l | sed '/\./ s/\.\{0,1\}0\{1,\}$//') != 0 ]; then
        echo "package_path ${package_path} not exist"
        exit__80_v0 1;
        __AF_exit80_v0__41_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__41_9" > /dev/null 2>&1
fi
    # 创建 SERVICE
    create_service_file__100_v0 "${install_dir}";
    __AF_create_service_file100_v0__45_5="$__AF_create_service_file100_v0";
    echo "$__AF_create_service_file100_v0__45_5" > /dev/null 2>&1
    # 创建安装目录
    create_dir_if_not_exists__99_v0 "${install_dir}";
    __AF_create_dir_if_not_exists99_v0__48_5="$__AF_create_dir_if_not_exists99_v0";
    echo "$__AF_create_dir_if_not_exists99_v0__48_5" > /dev/null 2>&1
    # Tar 解压包到安装目录
     tar -zxvf ${package_path} -C ${install_dir} ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "tar -zxvf ${package_path} -C ${install_dir} faild"
        exit__80_v0 1;
        __AF_exit80_v0__54_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__54_9" > /dev/null 2>&1
fi
    echo "tar success"
    # sed 字符串替换
    sed -i  "s/#HOST#/${host}/g" ${install_dir}/config.json;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "set host faild"
        exit__80_v0 1;
        __AF_exit80_v0__62_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__62_9" > /dev/null 2>&1
fi
    echo "install success"
     systemctl restart sgridnode ;
    __AS=$?;
if [ $__AS != 0 ]; then
        echo "systemctl restart sgridnode failed"
        exit__80_v0 1;
        __AF_exit80_v0__69_9="$__AF_exit80_v0";
        echo "$__AF_exit80_v0__69_9" > /dev/null 2>&1
fi
    echo "start sgridnode success"
