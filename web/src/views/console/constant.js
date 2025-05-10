export function  getServerType(type){
    if(type == 1){
        return "Node"
    }
    if(type == 2){
        return "Java"
    }
    if(type == 3){
        return "Binary"
    }
  }

export function getServerStatus(status){
    if(status == 1){
        return "Running"
    }
    if(status == 2){
        return "Stopped"
    }
    if(status == 3){
        return "Deleted"
    }
}

export function getServerNodeStatusType(status){
    if(status == 1){
        return "Success"
    }
    if(status == 2){
        return "Failed"
    }
    if(status == 3){
        return "Info"
    }
    if(status == 4){
        return "Warning"
    }
    if(status == 5){
        return "Patch"
    }
    if(status == 6){
        return "Check"
    }
}