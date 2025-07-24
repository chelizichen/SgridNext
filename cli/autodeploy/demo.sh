# 使用命令行参数
./sgridnext --DEPLOY_PATH=http://localhost:8080 --SERVER_ID=4 --SERVER_NAME=SgridTestJavaServer --PACKAGE_PATH=/path/to/SgridTestJavaServer.tar.gz

# 使用配置文件（原有方式）
./sgridnext

# 混合使用（命令行参数会覆盖配置文件）
./sgridnext --SERVER_ID=5 --PACKAGE_PATH=/new/path/package.tar.gz

# 查看帮助
./sgridnext -h