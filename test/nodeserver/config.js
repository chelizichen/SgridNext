const path = require('path')
const fs = require('fs')

function loadConfig(){
    const confDir = process.env.SGRID_CONF_DIR || process.cwd()
    const configPath = path.resolve(confDir, 'config.json')
    const configContent = fs.readFileSync(configPath, 'utf8')
    if (!configContent) throw new Error('No config.json found in ' + configPath)
    // parse config content
    const conf = JSON.parse(configContent)
    if (!conf) throw new Error('parse config error')
    return conf
}

module.exports = loadConfig