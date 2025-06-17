/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_327047008")

  // update collection data
  unmarshal({
    "indexes": [
      "CREATE INDEX `id_artist` ON `tracks` (`artists`)",
      "CREATE INDEX `id_album` ON `tracks` (`album`)"
    ]
  }, collection)

  // update field
  collection.fields.addAt(3, new Field({
    "cascadeDelete": true,
    "collectionId": "pbc_910879689",
    "hidden": false,
    "id": "relation22648455",
    "maxSelect": 999,
    "minSelect": 0,
    "name": "artists",
    "presentable": false,
    "required": true,
    "system": false,
    "type": "relation"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_327047008")

  // update collection data
  unmarshal({
    "indexes": [
      "CREATE INDEX `id_artist` ON `tracks` (`artist`)",
      "CREATE INDEX `id_album` ON `tracks` (`album`)"
    ]
  }, collection)

  // update field
  collection.fields.addAt(3, new Field({
    "cascadeDelete": true,
    "collectionId": "pbc_910879689",
    "hidden": false,
    "id": "relation22648455",
    "maxSelect": 999,
    "minSelect": 0,
    "name": "artist",
    "presentable": false,
    "required": true,
    "system": false,
    "type": "relation"
  }))

  return app.save(collection)
})
