https://github.com/golang/mock
mockgen -package mockmongo -destination mongo/mock/mongo.go unit-tests/mongo MongoStorage
mockgen -package mockothers -destination others/mock/others.go unit-tests/others OtherStorage