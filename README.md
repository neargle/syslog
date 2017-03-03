    
    info.GetSysLog("windows", "all")
    >> [map[time:03/02-17:16 username:asdasd hostname:aaaaaa status:false ip:-] map[hostname:111.111.11.111 status:false ip:- time:02/23-10:44 username:23333] ....

    usage:
    func GetSysLog(system string, starttime string) []map[string]string {
        // system : windows or linux
        // starttime : all or MM/dd-hh:ss
    }
    

