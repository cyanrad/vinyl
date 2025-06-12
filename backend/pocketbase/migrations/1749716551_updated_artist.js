/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_910879689")

  // update collection data
  unmarshal({
    "name": "artists"
  }, collection)

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_910879689")

  // update collection data
  unmarshal({
    "name": "artist"
  }, collection)

  return app.save(collection)
})
