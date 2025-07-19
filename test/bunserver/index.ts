import createServer from "./src/serve";

const server = createServer();

console.log(`Server is running on ${server.url}`);