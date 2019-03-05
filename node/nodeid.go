package node

import (
	"encoding/base64"

	"github.com/MagicalTux/seidan/db"
)

func NodeId() string {
	// get a node id from db
	id, err := db.SimpleGet([]byte("local"), []byte("node_id"))

	if err != nil {
		// most likely os.ErrNotExist
		id = GetKeyHash()                                    // will generate a key if needed
		db.SimpleSet([]byte("local"), []byte("node_id"), id) // we store the id so it's fast and doesn't require loading private key
	}

	return base64.RawURLEncoding.EncodeToString(id)
}
