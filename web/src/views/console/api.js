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
export function upsertConfig(data) {
  return request({ url: "/server/upsertConfig", data });
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
export function checkServerNodesStatus(data) {
  return request({ url: "/server/checkServerNodesStatus", data });
}
export function getServerNodesLog(data) {
  return request({ url: "/server/getServerNodesLog", data });
}
export function getServerNodes(data) {
  return request({ url: "/server/getServerNodes", data });
}
export function getServerConfigList(data) {
  return request({ url: "/server/getServerConfigList", data });
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

export function getConfigContent(data) {
  return request({ url: "/server/getConfigContent", data });
}

export function login(data) {
  return request({ url: "/login", data });
}




// cgroup

export function setCpuLimit(data) {
  return request({ url: "/server/cgroup/setCpuLimit", data });
}

export function getStatus(data) {
  return request({ url: "/server/cgroup/getStatus", data });
}

export function setMemoryLimit(data) {
  return request({ url: "/server/cgroup/setMemLimit", data });
}

// {
//   "data": {
//       "cpu": {
//           "usage": 1767523414000,
//           "usagePerSec": 1767.523414,
//           "shares": 0,
//           "throttled": 50249141
//       },
//       "memory": {
//           "usage": 331776,
//           "limit": 0,
//           "cache": 4096,
//           "swapUsage": 0,
//           "swapLimit": 0,
//           "oomEvents": 0
//       },
//       "io": {
//           "readBytes": 0,
//           "writeBytes": 0
//       },
//       "pids": {
//           "current": 5,
//           "limit": 18446744073709551615
//       },
//       "version": "v2",
//       "time": "2025-05-01T12:27:17.195937734+08:00"
//   },
//   "success": true
// }