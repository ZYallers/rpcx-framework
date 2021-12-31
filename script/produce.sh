#!/bin/bash

jsonFile="$(pwd)/service.json"
serviceName=""
serviceAddr=""
logDir=""

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

getJsonValue(){
    echo `cat ${jsonFile} | grep \"$1\"\: | sed -n '1p' | awk -F '": ' '{print $2}' | sed 's/,//g' | sed 's/"//g' | sed 's/ //g'`
}

helpFun(){
    echoFun "Operation:" title
    echoFun "    status                                  View service status" tip
    echoFun "    sync                                    Synchronization service vendor resources" tip
    echoFun "    build                                   Compile and generate service program" tip
    echoFun "    reload                                  Smooth restart service" tip
    echoFun "    quit                                    Stop service" tip
    echoFun "    help                                    View help information for the help command" tip
    echoFun "For more information about an action, use the help command to view it" tip
}

initFun(){
    echoFun "service config:" title
    if [[ ! -f "$jsonFile" ]];then
        echoFun "file [$jsonFile] is not exist" err
        exit 1
    fi

    serviceName=$(getJsonValue "name")
    if [[ "$serviceName" == "" ]];then
        echoFun "serviceName is empty" err
        exit 1
    fi
    echoFun "ServiceName: $serviceName" tip

    sip=$(ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:"|head -1)
    if [[ "$sip" == "" ]];then
        echoFun "SystemIP is empty" err
        exit 1
    fi

    serviceAddr=$(getJsonValue "addr")
    if [[ "$serviceAddr" == "" ]];then
        echoFun "ServiceAddr is empty" err
        exit 1
    fi

    port=`echo ${serviceAddr}|awk -F ':' '{print $2}'`
    serviceAddr="$sip:$port"
    echoFun "ServiceAddr: $serviceAddr" tip

    logDir=$(getJsonValue "logDir")
    if [[ "$logDir" == "" ]];then
        echoFun "LogDir is empty" err
        exit 1
    fi

    echoFun "LogDir: $logDir" tip
}

statusFun(){
    echoFun "ps process:" title
    if [[ `pgrep ${serviceName}|wc -l` -gt 0 ]];then
        ps -p $(pgrep ${serviceName}|sed ':t;N;s/\n/,/;b t'|sed -n '1h;1!H;${g;s/\n/,/g;p;}') -o user,pid,ppid,%cpu,%mem,vsz,rss,tty,stat,start,time,command
    fi

    echoFun "lsof process:" title
    port=`echoFun ${serviceAddr}|awk -F ':' '{print $2}'`
    lsof -i:${port}
}

syncFun(){
    echoFun "go mod vendor:" title
    if [[ ! -f "./go.mod" ]];then
        echoFun "go.mod is not exist" err
        exit 1
    fi
    go mod tidy
    go mod vendor
    echoFun "go mod vendor finished" ok
}

buildFun(){
    echoFun "git pull:" title
    branch=$1
    env=$2

    if [[ "$branch" == "" ]];then
        echoFun "branch of the build is empty" err
        exit 1
    fi
    if [[ "$branch" == "local" ]];then
        echoFun "ignore git pull, direct build by local" tip
    else
        git remote update origin --prune # Update remote branch list
        git checkout ${branch} # Switch branch
        git pull # Pull the latest version
        echoFun "git pull [$branch] finish" ok
    fi

    echoFun "build runner:" title
    tmpName="${serviceName}_$(date +'%Y%m%d%H%M%S')"
    if [[ "$env" == "dev" ]];then
        echoFun 'build in develop environment' tip
        CGO_ENABLED=0 go build -v -installsuffix cgo -ldflags '-w' -i -o ./bin/${tmpName} -tags=jsoniter ./main.go
    else
        # Build compilation parameter reference:
        # Dependency free compilation：https://blog.csdn.net/weixin_42506905/article/details/93135684
        # Detailed explanation of build parameters：https://blog.csdn.net/zl1zl2zl3/article/details/83374131
        # Ldflags parameter：https://blog.csdn.net/javaxflinux/article/details/89177863
        CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-w' -i -o ./bin/${tmpName} -tags=jsoniter ./main.go
    fi

    if [[ ! -f "./bin/${tmpName}" ]];then
        echoFun "build tmp runner ($(pwd)/bin/${tmpName}) failed" err
        exit 1
    fi

    mv -f ./bin/${tmpName} ./bin/${serviceName}
    if [[ ! -f "./bin/${serviceName}" ]];then
        echoFun "mv tmp runner failed" err
        exit 1
    fi

    echoFun "build runner ($(pwd)/bin/${serviceName}) finished" ok
}

sendMsg(){
    envVar="$(echo $hxsenv)"
    if [[ "$envVar" == "" ]];then
        envVar="development"
    fi
    env="Env: $envVar"
    name="Name: $serviceName"
    addr="Addr: $serviceAddr"
    hostnameVar=$(hostname)
    if [[ "$hostnameVar" == "" ]];then
        envVar="unknown"
    fi
    hostname="HostName: $hostnameVar"
    sipVar=$(ifconfig -a|grep inet|grep -v 127.0.0.1|grep -v inet6|awk '{print $2}'|tr -d "addr:"|head -1)
    if [[ "$sipVar" == "" ]];then
        sipVar="unknown"
    fi
    sip="SystemIP: $sipVar"
    time="Time: $(date "+%Y/%m/%d %H:%M:%S")"
    token=$(getJsonValue "gracefulRobotToken")
    url="https://oapi.dingtalk.com/robot/send?access_token=$token"
    content="$1\n---------------------------\n$env\n$name\n$addr\n$hostname\n$time\n$sip"
    cnt=$(echo ${content//\"/\\\"})
    header="Content-Type: application/json"
    curl -o /dev/null -m 3 -s "$url" -H "$header" -d "{\"msgtype\":\"text\",\"text\":{\"content\":\"$cnt\"}}"
}

reloadFun(){
    echoFun "runner reloading:" title

    if [[ ! -f "./bin/$serviceName" ]];then
        echoFun "runner [`pwd`/bin/$serviceName] is not exist" err
        exit 1
    fi

    if [[ ! -d "$logDir" ]];then
        mkdir -p ${logDir}
    fi
    if [[ ! -d "$logDir" ]];then
        echoFun "logDir [$logDir] is not exist" err
        exit 1
    fi

    logfile=${logDir}/${serviceName}.log
    if [[ ! -f "$logfile" ]];then
        touch ${logfile}
    fi
    echoFun "logfile: $logfile" tip

    if [[ ! -x "./bin/$serviceName" ]];then
        chmod u+x ./bin/${serviceName}
    fi

    quitFun

    # Prevent Jenkins from killing all derived processes after the end of build by default
    export BUILD_ID=dontKillMe

    nohup ./bin/${serviceName} >> ${logfile} 2>&1 &
    echoFun "service $serviceName($serviceAddr) is reloaded, pid: `echo $!`" ok

    # Check whether the health interface is accessed normally
    sleep 3s
    id=$(date +%Y%m%d%H%M%S)
    post="{\"jsonrpc\":\"2.0\",\"method\":\"$serviceName.health\",\"params\":{},\"id\":$id}"
    resp=`curl -m 3 -s "http://$serviceAddr" --header "X-JSONRPC-2.0:true" -d "$post"`
    echoFun "curl \"http://$serviceAddr\" health: $resp" tip
    sendMsg "curl \"http://$serviceAddr\" health: $resp"
}

quitFun(){
    port=`echo ${serviceAddr}|awk -F ':' '{print $2}'`
    counter=0
    while true;
    do
        pid=`lsof -i tcp:${port}|grep LISTEN|awk '{print $2}'`
        if [[ ${pid} -gt 0 ]];then
            if [[ ${counter} -ge 30 ]];then
                kill -9 ${pid}
                echoFun "service $serviceName has been killed for 30s and is ready to be forcibly killed" tip
                sendMsg "service $serviceName has been killed for 30s and is ready to be forcibly killed"
                break
            else
                kill ${pid}
                counter=$(($counter+1))
                echoFun "killing service $serviceName($port), pid($pid), $counter tried" tip
                sleep 1s
            fi
        else
            echoFun "service $serviceName($port) service is stopped" ok
            break
        fi
    done
}

initFun
case $1 in
    status)
        statusFun
    ;;
    sync)
        syncFun
    ;;
    build)
        buildFun $2 $3
    ;;
    quit)
        quitFun
    ;;
    reload)
        reloadFun
    ;;
    *)
        helpFun
    ;;
esac
