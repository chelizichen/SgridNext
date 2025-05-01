import axios from "axios";

const request = axios.create({
  baseURL: "/api",
  timeout: 60 * 1000,
  method: "post",
});

request.interceptors.response.use(
  (data) => {
    return data.data;
  },
  (error) => {
    return Promise.reject(error);
  }
);

export function createServer(data) {
  return request({ url: "/server/createServer", data });
}
export function uploadPackage(data) {
  return request({ url: "/server/uploadPackage", data });
}
export function createConfig(data) {
  return request({ url: "/server/createConfig", data });
}
export function createServerNode(data) {
  return request({ url: "/server/createServerNode", data });
}
export function createGroup(data) {
  return request({ url: "/server/createGroup", data });
}
export function createNode(data) {
  return request({ url: "/server/createNode", data });
}
export function deployServer(data) {
  return request({ url: "/server/deployServer", data });
}
export function stopServer(data) {
  return request({ url: "/server/stopServer", data });
}
export function restartServer(data) {
  return request({ url: "/server/restartServer", data });
}
export function getServerNodesStatus(data) {
  return request({ url: "/server/getServerNodesStatus", data });
}
export function getServerNodesLog(data) {
  return request({ url: "/server/getServerNodesLog", data });
}
export function getServerNodes(data) {
  return request({ url: "/server/getServerNodes", data });
}
export function getServerConfigList(data) {
  return request({ url: "/server/getServList", data });
}
export function getServerPackageList(data) {
  return request({ url: "/server/getServerPackageList", data });
}
export function getServerList(data) {
  return request({ url: "/server/getServerList", data });
}
export function getNodeList(data) {
  return request({ url: "/server/getNodeList", data });
}
export function getNodeLoadDetail(data) {
  return request({ url: "/server/getNodeLoadDetail", data });
}
export function updateConfig(data) {
  return request({ url: "/server/updateConfig", data });
}
export function getGroupList(data) {
  return request({ url: "/server/getGroupList", data });
}
export function getServerInfo(data) {
  return request({ url: "/server/getServerInfo", data });
}
export function login(data) {
  return request({ url: "/login", data });
}


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

// cgroup

export function setCpuLimit(data) {
  return request({ url: "/server/cgroup/setCpuLimit", data });
}