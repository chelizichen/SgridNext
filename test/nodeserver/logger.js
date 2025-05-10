const { createLogger, format } = require("winston");
const DailyRotateFile = require("winston-daily-rotate-file");
const path = require("path");
const initLogger = (logName) => {
  let logDir = process.env.SGRID_LOG_DIR || path.resolve(process.cwd(), "log");
  const filePath = path.resolve(logDir, `${logName}.log`);
  return createLogger({
    level: "info", // 设置日志级别
    format: format.combine(
      format.timestamp({
        format: "YYYY-MM-DD HH:mm:ss",
      }),
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
    ],
  });
};

// SERVICE LOGGER
const logger = {};
logger.data = initLogger("data");
module.exports = logger;