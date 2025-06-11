/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_327047008")

  // update collection data
  unmarshal({
    "indexes": [
      "CREATE INDEX `id_artist` ON `tracks` (`artist`)",
      "CREATE INDEX `id_album` ON `tracks` (`album`)"
    ]
  }, collection)

  // add field
  collection.fields.addAt(5, new Field({
    "hidden": false,
    "id": "file2366146245",
    "maxSelect": 1,
    "maxSize": 0,
    "mimeTypes": [],
    "name": "cover",
    "presentable": false,
    "protected": false,
    "required": false,
    "system": false,
    "thumbs": [],
    "type": "file"
  }))

  // add field
  collection.fields.addAt(7, new Field({
    "cascadeDelete": true,
    "collectionId": "pbc_3287366145",
    "hidden": false,
    "id": "relation966291011",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "album",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "relation"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_327047008")

  // update collection data
  unmarshal({
    "indexes": [
      "CREATE INDEX `id_artist` ON `tracks` (`artist`)"
    ]
  }, collection)

  // remove field
  collection.fields.removeById("file2366146245")

  // remove field
  collection.fields.removeById("relation966291011")

  return app.save(collection)
})
