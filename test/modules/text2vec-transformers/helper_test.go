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
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate/client/objects"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/test/helper"
)

func setupClient(uri string) {
	host, port := "", ""
	res := strings.Split(uri, ":")
	if len(res) == 2 {
		host, port = res[0], res[1]
	}
	helper.ServerHost = host
	helper.ServerPort = port
}

func createObjectClass(t *testing.T, class *models.Class) {
	helper.CreateObjectClass(t, class)
}

func createObject(t *testing.T, object *models.Object) {
	helper.CreateObject(t, object)
}

func deleteObjectClass(t *testing.T, class string) {
	helper.DeleteObjectClass(t, class)
}

func assertGetObjectEventually(t *testing.T, uuid strfmt.UUID) *models.Object {
	var (
		resp *objects.ObjectsGetOK
		err  error
	)

	checkThunk := func() interface{} {
		resp, err = helper.Client(t).Objects.ObjectsGet(objects.NewObjectsGetParams().WithID(uuid), nil)
		return err == nil
	}

	helper.AssertEventuallyEqual(t, true, checkThunk)

	var object *models.Object

	helper.AssertRequestOk(t, resp, err, func() {
		object = resp.Payload
	})

	return object
}
