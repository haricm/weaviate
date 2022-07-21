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

package docker

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
)

const (
	// TEST_WEAVIATE_IMAGE can be passed to tests to spin up docker compose with given image
	TEST_WEAVIATE_IMAGE = "TEST_WEAVIATE_IMAGE"
	// TEST_TEXT2VEC_TRANSFORMERS_IMAGE adds ability to pass a custom image to module tests
	TEST_TEXT2VEC_TRANSFORMERS_IMAGE = "TEST_TEXT2VEC_TRANSFORMERS_IMAGE"
)

type Compose struct {
	enableModules           []string
	defaultVectorizerModule string
	withMinIO               bool
	withTransformers        bool
	withWeaviate            bool
}

func New() *Compose {
	return &Compose{enableModules: []string{}}
}

func (d *Compose) WithMinIO() *Compose {
	d.withMinIO = true
	return d
}

func (d *Compose) WithText2VecTransformers() *Compose {
	d.withTransformers = true
	d.enableModules = append(d.enableModules, Text2VecTransformers)
	d.defaultVectorizerModule = Text2VecTransformers
	return d
}

func (d *Compose) WithWeaviate() *Compose {
	d.withWeaviate = true
	return d
}

func (d *Compose) Start(ctx context.Context) (*DockerCompose, error) {
	networkName := "weaviate-module-acceptance-tests"
	network, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:     networkName,
			Internal: false,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "network: %s", networkName)
	}
	envSettings := make(map[string]string)
	containers := []*DockerContainer{}
	if d.withMinIO {
		container, err := startMinIO(ctx, networkName)
		if err != nil {
			return nil, errors.Wrapf(err, "start %s", MinIO)
		}
		containers = append(containers, container)
	}
	if d.withTransformers {
		image := os.Getenv(TEST_TEXT2VEC_TRANSFORMERS_IMAGE)
		container, err := startT2VTransformers(ctx, networkName, image)
		if err != nil {
			return nil, errors.Wrapf(err, "start %s", Text2VecTransformers)
		}
		for k, v := range container.envSettings {
			envSettings[k] = v
		}
		containers = append(containers, container)
	}
	if d.withWeaviate {
		image := os.Getenv(TEST_WEAVIATE_IMAGE)
		container, err := startWeaviate(ctx, d.enableModules, d.defaultVectorizerModule,
			envSettings, networkName, image)
		if err != nil {
			return nil, errors.Wrapf(err, "start %s", Weaviate)
		}
		containers = append(containers, container)
	}

	return &DockerCompose{network, containers}, nil
}
