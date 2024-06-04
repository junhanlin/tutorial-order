package component

import (
	"github.com/bwmarrin/snowflake"
	log "github.com/sirupsen/logrus"
)

func NewSnowflake() *snowflake.Node {
	// Setup snowflake node
	snowflakeNode, err := snowflake.NewNode(1)
	if err != nil {
		log.WithError(err).Fatal("Error creating snowflake node")
	}
	return snowflakeNode
}
