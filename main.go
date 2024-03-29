/*
 * Copyright 2020 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"flag"
	"github.com/SENERGY-Platform/marshaller/lib"
	"github.com/SENERGY-Platform/marshaller/lib/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configLocation := flag.String("config", "config.json", "configuration file")
	flag.Parse()

	conf, err := config.Load(*configLocation)
	if err != nil {
		log.Fatal("ERROR: unable to load config", err)
	}

	ctx, shutdown := context.WithCancel(context.Background())

	closed, err := lib.Start(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}

	err = lib.StartCacheInvalidator(ctx, conf)
	if err != nil {
		log.Println("WARNING: unable to start cache invalidator", err)
	}

	go func() {
		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		sig := <-shutdownSignal
		log.Println("received shutdown signal", sig)
		shutdown()
	}()

	<-closed.Done()
	log.Println("server closed")

}
