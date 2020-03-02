// package infinote holds platform wide source code

//go:generate ../bin/go-bindata -prefix migrations/ -pkg bindata -nocompress -o ./bindata/bindata.go migrations
//go:generate ../bin/sqlboiler ../bin/sqlboiler-psql --wipe --tag db --config ./sqlboiler.toml --output ./db
//go:generate go run github.com/vektra/mockery/cmd/mockery -name UserStorer
//go:generate go run github.com/vektra/mockery/cmd/mockery -name TokenStorer
//go:generate go run github.com/vektra/mockery/cmd/mockery -name BlacklistProvider
//go:generate go run github.com/vektra/mockery/cmd/mockery -name AuthProvider

package infinote
