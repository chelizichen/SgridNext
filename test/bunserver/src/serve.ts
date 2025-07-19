import logger from "./components/logger/main";
import routes from "./domain/controller";

function createServer() {
    logger.data.info("PATH %s",Object.keys(routes));
    let PORT = process.env.SGRID_TARGET_PORT || 3000;
    let HOST = process.env.SGRID_TARGET_HOST || '0.0.0.0';
    const server = Bun.serve({
        port: PORT,
        hostname: HOST,
        routes:Object.assign({},routes)
    });

    return server;
}

export default createServer;