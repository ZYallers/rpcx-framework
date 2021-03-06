#!/bin/bash
# Initialize the script, pull the data, and generate the compilation and deployment script

echoFun(){
    str=$1
    color=$2
    case ${color} in
        ok)
            echo -e "\033[32m $str \033[0m"
        ;;
        err)
            echo -e "\033[31m $str \033[0m"
        ;;
        tip)
            echo -e "\033[33m $str \033[0m"
        ;;
        title)
            echo -e "\033[42;34m $str \033[0m"
        ;;
        *)
            echo "$str"
        ;;
    esac
}

echoFun "go version" title
go version

echoFun "go env" title
go env

echoFun "current path: $(pwd)" title
echoFun "download produce.sh" title
src='https://raw.githubusercontent.com/ZYallers/rpcx-framework/master/script/produce.sh'
des='./bin/produce.sh'
curl -o ${des} ${src}
if [[ ! -f "$des" ]];then
    echoFun "download produce.sh($(pwd)/$des) failed" err
    exit 1
fi
chmod u+x ${des}
echoFun "download produce.sh($(pwd)/$des) finished" ok
./bin/produce.sh help

