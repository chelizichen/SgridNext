import winston, { createLogger, format } from "winston";
import DailyRotateFile from "winston-daily-rotate-file";
import path from "path";
const initLogger = (logName: string) => {
    let filePath = "";
    if(process.env.SGRID_LOG_DIR){
        filePath = path.resolve(process.env.SGRID_LOG_DIR, `${logName}.log`);
    }else{
        filePath = path.resolve(process.cwd(), "log", `${logName}.log`);
    }
  return createLogger({
    level: "info", // 设置日志级别
    format: format.combine(
      format.timestamp({
        format: "YYYY-MM-DD HH:mm:ss",
      }),
      format.splat(),
      format.printf(
        (info) => `${info.timestamp} ${info.level}: ${info.message}`
      )
    ),
    transports: [
      new DailyRotateFile({
        filename: filePath.replace(/\.log$/, "-%DATE%.log"),
        datePattern: "YYYY-MM-DD",
        zippedArchive: true,
        maxSize: "20m",
        maxFiles: "14d",
      }),
      // 在开发环境下添加控制台输出
      ...(!process.env.SGRID_LOG_DIR ? [
        new winston.transports.Console({
          format: format.combine(
            format.colorize(),
            format.timestamp({
              format: "YYYY-MM-DD HH:mm:ss",
            }),
            format.splat(),
            format.printf(
              (info) => `${info.timestamp} ${info.level}: ${info.message}`
            )
          )
        })
      ] : [])
    ],
  });
};

const logger = {
  data: initLogger("data"),
  error: initLogger("error"),
};

export default logger;