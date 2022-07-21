//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2022 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package test

import (
	"context"
	"os"
	"testing"

	"github.com/semi-technologies/weaviate/test/docker"
	"github.com/sirupsen/logrus"
)

const minioEndpoint = "MINIO_ENDPOINT"

func TestMain(m *testing.M) {
	logger := logrus.New()
	if os.Getenv("LOG_FORMAT") != "text" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	ctx := context.Background()
	compose, err := docker.New().WithMinIO().Start(ctx)
	if err != nil {
		logger.WithError(err).Panic("cannot start")
	}

	os.Setenv(minioEndpoint, compose.GetMinIO().URI())
	code := m.Run()

	if err := compose.Terminate(ctx); err != nil {
		logger.WithError(err).Error("cannot terminate")
	}

	os.Exit(code)
}
