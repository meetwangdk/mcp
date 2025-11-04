package data

import (
	"fmt"
	"github.com/google/wire"
	"github.com/olivere/elastic"
	"hcs-agent/pkg/config"
	"hcs-agent/pkg/log"
)

// Set is data providers.
var Set = wire.NewSet(
	NewEsDb,
	NewDemoRepo,
)

type EsDb struct {
	db *elastic.Client
}

func NewEsDb(config *config.Config, log *log.Logger) *EsDb {
	esClient, err := elastic.NewClient(
		elastic.SetURL(config.Es.Endpoint), elastic.SetBasicAuth(config.Es.Username, config.Es.Password), elastic.SetSniff(false))
	if err != nil {
		log.Error(fmt.Sprintf("init es client error %s", err))
		panic(err)
	}
	return &EsDb{
		db: esClient,
	}
}
