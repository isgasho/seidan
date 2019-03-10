package node

type nodeStatus int

const (
	STATUS_OFFLINE nodeStatus = iota
	STATUS_INIT
	STATUS_CONNECT
	STATUS_FOLLOWER
	STATUS_CANDIDATE
	STATUS_LEADER
)

type Monitor struct {
	status nodeStatus
}
