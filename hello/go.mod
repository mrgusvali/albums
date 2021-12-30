module github.com/mrgusvali/hello

go 1.17

replace github.com/mrgusvali/rediscli => ../redis-cli

require github.com/mrgusvali/rediscli v0.0.0-00010101000000-000000000000

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.4 // indirect
)
