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
  },
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

export function updateServerNode(data) {
  return request({ url: "/server/updateServerNode", data });
}

export function updateServer(data) {
  return request({ url: "/server/updateServer", data });
}
export function deleteServerNode(data){
  return request({ url: "/server/deleteServerNode",data });
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

export function updateMachineNodeStatus(data){
  return request({ url: "/server/updateNode", data });
}

export function updateMachineNodeAlias(data){
  return request({ url: "/server/updateNodeAlias", data });
}

// downloadFile({
//     serverId:1,
//     fileName:"waterfull.log",
//     type:1,
//     host:"10.0.12.17"
// })
export function downloadFile(data) {
  return new Promise((resolve, reject) => {
    request({ url: "/server/downloadFile", data })
      .then((res) => {
        // 将res 下载成一个文件
        const blob = new Blob([res], { type: "application/octet-stream" });
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = data.fileName;
        a.click();
        URL.revokeObjectURL(url);
        resolve(true);
      })
      .catch((err) => {
        reject(err);
      });
  });
}

export function getFileList(data) {
  return request({ url: "/server/getFileList", data });
}

export function getLog(data) {
  return request({ url: "/server/getLog", data });
}

export function getSyncStatus(data){
  return request({ 
    url: "/server/getSyncStatus",
    data 
  });
}


export function syncUploadFile(data){
  return request({ 
    url: "/server/syncUploadFile",
    data 
  });
}

// 主控配置管理
export function getMainConfig() {
  return request({ url: "/config/getMainConfig" });
}

export function updateMainConfig(data) {
  return request({ url: "/config/updateMainConfig", data });
}

export function getConfigItem(data) {
  return request({ url: "/config/getConfigItem", data });
}

export function setConfigItem(data) {
  return request({ url: "/config/setConfigItem", data });
}

// 探针相关API
export function runProbeTask() {
  return request({ url: "/probe/runProbeTask" });
}

// 获取节点资源信息
export function getNodeResource(nodeId) {
  return request({ url: "/resource/getNodeResource", data: { nodeId } });
}

// ========== 文档管理 API ==========
export function uploadDocument(formData) {
  return axios.post("/api/document/upload", formData, {
    headers: { "Content-Type": "multipart/form-data" },
    timeout: 60 * 1000,
  }).then(res => res.data);
}

export function createDocument(data) {
  return request({ url: "/document/create", data });
}

export function updateDocument(data) {
  return request({ url: "/document/update", data });
}

export function deleteDocument(data) {
  return request({ url: "/document/delete", data });
}

export function getDocumentList() {
  return request({ url: "/document/list" });
}

export function getDocument(id) {
  return request({ url: "/document/get", data: { id } });
}

export function downloadDocument(id) {
  return axios.get(`/api/document/download?id=${id}`, {
    responseType: "blob",
  });
}

export function linkDocumentToServer(data) {
  return request({ url: "/document/link", data });
}

export function getDocumentServerRelations(documentId) {
  return axios.get(`/api/document/relations?documentId=${documentId}`).then(res => res.data);
}