/*
Sniperkit-Bot
- Status: analyzed
*/

package indexer

import (
	"fmt"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/index/upsidedown"
	"github.com/blevesearch/bleve/mapping"

	"github.com/sniperkit/snk.fork.gostore-contrib/indexer/badger"
	_ "github.com/sniperkit/snk.fork.gostore-contrib/indexer/badger"
)

// NewIndexer creates a new indexer
func NewBadgerIndexer(indexPath string) *Indexer {
	indexMapping := bleve.NewIndexMapping()
	return NewBadgerIndexerWithMapping(indexPath, indexMapping)
}

// NewIndexer creates a new indexer
func NewBadgerIndexerWithMapping(indexPath string, indexMapping mapping.IndexMapping) *Indexer {
	index, err := bleve.Open(indexPath)
	if err != nil {
		logger.Debug("Error opening indexpath", "path", indexPath, "verbose", string(err.Error()))
		if err == bleve.ErrorIndexMetaMissing || err == bleve.ErrorIndexPathDoesNotExist {
			logger.Debug(fmt.Sprintf("Creating new badger index at %s ...", indexPath))
			// indexMapping.DefaultAnalyzer = "keyword"
			kvconfig := map[string]interface{}{}

			index, err = bleve.NewUsing(indexPath, indexMapping, upsidedown.Name, badger.Name, kvconfig)

			if err != nil {
				logger.Warn("Index could not be created", "path", indexPath, "err", string(err.Error()))
				if err != bleve.ErrorIndexPathExists {
					panic(err)
				}
				return nil
			}

		} else {
			panic(err)
		}
	}
	logger.Debug("opening existing index", "stats", index.Stats())
	return &Indexer{index: index}
}
