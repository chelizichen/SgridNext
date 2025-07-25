package logger

var CMD = CreateLogger("cmd")
var PatchUtils = CreateLogger("patchutils")
var App = CreateLogger("app")
var Server = CreateLogger("server")
var Mapper = CreateLogger("mapper")
var Cgroup = CreateLogger("cgroup")
var Hook_Cgroup = CreateLogger("hook_cgroup")
var Config = CreateLogger("config")

var Alive = CreateLogger("alive")

var RPC = CreateLogger("rpc")

var Docker = CreateLogger("docker")