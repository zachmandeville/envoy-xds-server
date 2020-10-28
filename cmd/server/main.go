//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"context"
	"flag"

	serverv2 "github.com/envoyproxy/go-control-plane/pkg/server/v2"

	v2Cache "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	testv3 "github.com/envoyproxy/go-control-plane/pkg/test/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stevesloka/envoy-xds-server/internal/processor"
	"github.com/stevesloka/envoy-xds-server/internal/server"
	"github.com/stevesloka/envoy-xds-server/internal/watcher"
)

var (
	l log.FieldLogger

	watchDirectoryFileName string
	port                   uint
	basePort               uint
	mode                   string

	nodeID string
)

func init() {
	l = log.New()
	log.SetLevel(log.DebugLevel)

	// The port that this xDS server listens on
	flag.UintVar(&port, "port", 9002, "xDS management server port")

	// Tell Envoy to use this Node ID
	flag.StringVar(&nodeID, "nodeID", "test-id", "Node ID")

	// Define the directory to watch for Envoy configuration files
	flag.StringVar(&watchDirectoryFileName, "watchDirectoryFileName", "config/config.yaml", "full path to directory to watch for files")
}

func main() {
	flag.Parse()

	// Create a cachev3
	cachev3 := cache.NewSnapshotCache(false, cache.IDHash{}, l)
	cachev2 := v2Cache.NewSnapshotCache(false, v2Cache.IDHash{}, l)

	// Create a processor
	proc := processor.NewProcessor(
		cachev3, cachev2, nodeID, log.WithField("context", "processor"))

	// Create initial snapshot from file
	proc.ProcessFile(watcher.NotifyMessage{
		Operation: watcher.Create,
		FilePath:  watchDirectoryFileName,
	})

	// Notify channel for file system events
	notifyCh := make(chan watcher.NotifyMessage)

	go func() {
		// Watch for file changes
		watcher.Watch(watchDirectoryFileName, notifyCh)
	}()

	go func() {
		// Run the xDS server
		ctx := context.Background()
		cb := &testv3.Callbacks{Debug: true}
		srv := serverv3.NewServer(ctx, cachev3, cb)
		srvv2 := serverv2.NewServer(ctx, cachev2, nil)
		server.RunServer(ctx, srv, srvv2, port)
	}()

	for {
		select {
		case msg := <-notifyCh:
			proc.ProcessFile(msg)
		}
	}
}
