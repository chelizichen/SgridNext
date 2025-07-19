import logger from "../../components/logger/main";

const demoData = {
    bugs: [
        {
          id: 1,
          title: "SgridNext 发布包后内存泄露",
          description: "SgridNext 发布包后内存泄露",
          status: "open",
          priority: "high",
          createdAt: "2021-01-01",
          updatedAt: "2021-01-01",
        },
        {
          id: 2,
          title: "SgridNode 内存泄露",
          description: "SgridNode 内存泄露",
          status: "open",
          priority: "high",
          createdAt: "2021-01-01",
          updatedAt: "2021-01-01",
        },
        {
          id: 3,
          title: "SgridNext Web 缺少更新服务信息界面",
          description: "SgridNext Web 缺少更新服务信息界面",
          status: "open",
          priority: "high",
          createdAt: "2021-01-01",
          updatedAt: "2021-01-01",
        },
        {
          id: 4,
          title: "SgridNext 发布包失败、上传包失败",
          description: "SgridNext 发布包失败、上传包失败",
          status: "open",
          priority: "high",
          createdAt: "2021-01-01",
          updatedAt: "2021-01-01",
        },
      ],
      total:4,
}

const bugService = {
  getBugs() {
    return {
        ...demoData
    };
  },
  fix(id:number){
    logger.data.info("fix bug %s",id);
    return true
  },
  close(id:number){
    logger.data.info("close bug %s",id);
    return true
  }
};

export default bugService;
