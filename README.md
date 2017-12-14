    
    info.GetSysLog("windows", "all")
    >> [map[time:03/02-17:16 username:asdasd hostname:aaaaaa status:false ip:-] map[hostname:111.111.11.111 status:false ip:- time:02/23-10:44 username:23333] ....

    usage:
    func GetSysLog(system string, starttime string) []map[string]string {
        // system : windows or linux
        // starttime : all or MM/dd-hh:ss
    }
    

Windows调用了Powershell, Linux目前支持CentOS，其他应该也适用，不过没测试...
我看到有朋友关注了该项目说一下，这个项目已经换了其他的实现方法，重构了一遍，会作为某个HIDS项目的一部分，近期会开源。名为 "驭龙HIDS"，感兴趣的朋友可以关注一下。
