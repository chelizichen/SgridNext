package state

import "sync/atomic"

var NodeServerState = atomic.Int32{}

const NODE_STATE_STAYBY = 0
const NODE_STATE_ONLINE = 1
const NODE_STATE_DONE_INIT = 2
