
const express = require('express')
const app = express()
const loadConfig = require('./config')
const logger = require('./logger')
const port = process.env.SGRID_TARGET_PORT || 3000
const host = process.env.SGRID_TARGET_HOST || '0.0.0.0'

const conf = loadConfig()
logger.data.info('conf',conf);


app.get('/', (req, res) => {
  res.send('Hello World!')
})

function fib(n) {
  if (n === 0) return 0
  if (n === 1) return 1
  return fib(n - 1) + fib(n - 2)
}

app.get('/fib', (req, res) => {
  let n = req.query.n || 10
  let result = fib(n)
  res.send(result)
})


app.get('/bigmemmory', (req, res) => {
  let n = req.query.n || 10
  let result = fib(n)
  let bigmemory = new Array(1000000000)
  bigmemory.fill(1)
  result = result + bigmemory.length
  res.send(result)  
})

logger.data.info('host %s',host);
logger.data.info('port %s',port);

app.listen(port,host, () => {
  logger.data.info(`Example app listening on bind address ${host}:${port}`)
})

process.on("uncaughtException", (err) => {
  logger.data.error("uncaughtException", err);
});
process.on("unhandledRejection", (reason, promise) => {
  logger.data.error("unhandledRejection", reason);
})