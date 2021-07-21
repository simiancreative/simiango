let res = [
  db.products.drop(),
  db.products.createIndex({ name: 1 }, { unique: true }),
  db.products.createIndex({ id: 1 }),
  db.products.insert({ id: 1, name: 'Car' }),
  db.products.insert({ id: 2, name: 'Truck' }),
  db.products.insert({ id: 3, name: 'Motorcycle' }),
  db.products.insert({ id: 4, name: 'Bicycle' }),
  db.products.insert({ id: 5, name: 'Horse' }),
  db.products.insert({ id: 6, name: 'Boat' }),
  db.products.insert({ id: 7, name: 'Plane' }),
  db.products.insert({ id: 8, name: 'Scooter' }),
  db.products.insert({ id: 9, name: 'Skateboard' }),
]

printjson(res)