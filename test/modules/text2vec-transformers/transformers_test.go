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
	"os"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/test/helper"
	graphqlhelper "github.com/semi-technologies/weaviate/test/helper/graphql"
	"github.com/stretchr/testify/assert"
)

const (
	dune                  strfmt.UUID = "67b79643-cf8b-4b22-b206-6e63dbb4e000"
	projectHailMary       strfmt.UUID = "67b79643-cf8b-4b22-b206-6e63dbb4e001"
	theLordOfTheIceGarden strfmt.UUID = "67b79643-cf8b-4b22-b206-6e63dbb4e002"
)

func Test_T2VTransformers(t *testing.T) {
	setupClient(os.Getenv(weaviateEndpoint))
	createObjectClass(t, &models.Class{
		Class: "Books",
		ModuleConfig: map[string]interface{}{
			"text2vec-transformers": map[string]interface{}{
				"vectorizeClassName": true,
			},
		},
		Properties: []*models.Property{
			{
				Name:     "title",
				DataType: []string{"string"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": false,
					},
				},
			},
			{
				Name:     "description",
				DataType: []string{"string"},
				ModuleConfig: map[string]interface{}{
					"text2vec-transformers": map[string]interface{}{
						"skip": false,
					},
				},
			},
		},
	})
	defer deleteObjectClass(t, "Books")

	t.Run("add data to Books schema", func(t *testing.T) {
		createObject(t, &models.Object{
			Class: "Books",
			ID:    dune,
			Properties: map[string]interface{}{
				"title":       "Dune",
				"description": "Dune is a 1965 epic science fiction novel by American author Frank Herbert.",
			},
		})
		createObject(t, &models.Object{
			Class: "Books",
			ID:    projectHailMary,
			Properties: map[string]interface{}{
				"title":       "Project Hail Mary",
				"description": "Project Hail Mary is a 2021 science fiction novel by American novelist Andy Weir.",
			},
		})
		createObject(t, &models.Object{
			Class: "Books",
			ID:    theLordOfTheIceGarden,
			Properties: map[string]interface{}{
				"title":       "The Lord of the Ice Garden",
				"description": "The Lord of the Ice Garden (Polish: Pan Lodowego Ogrodu) is a four-volume science fiction and fantasy novel by Polish writer Jaroslaw Grzedowicz.",
			},
		})

		assertGetObjectEventually(t, dune)
		assertGetObjectEventually(t, projectHailMary)
		assertGetObjectEventually(t, theLordOfTheIceGarden)
	})

	t.Run("query Books data with nearText", func(t *testing.T) {
		result := graphqlhelper.AssertGraphQL(t, helper.RootAuth, `
			{
				Get {
					Books(
						nearText: {
							concepts: ["Andy Weir"]
							distance: 0.3
						}
					){
						title
					}
				}
			}
		`)
		books := result.Get("Get", "Books").AsSlice()

		expected := []interface{}{
			map[string]interface{}{"title": "Project Hail Mary"},
		}

		assert.ElementsMatch(t, expected, books)
	})
}
