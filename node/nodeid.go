package node

import (
	"github.com/MagicalTux/seidan/db"
	"github.com/google/uuid"
)

func NodeId() string {
	// get a node id from db
	v, err := db.SimpleGet([]byte("local"), []byte("node_id"))
	var id uuid.UUID

	if err != nil {
		// most likely os.ErrNotExist
		id = uuid.Must(uuid.NewRandom())
		db.SimpleSet([]byte("local"), []byte("node_id"), id[:])
		return id.String()
	}

	copy(id[:], v)
	return id.String()
}
