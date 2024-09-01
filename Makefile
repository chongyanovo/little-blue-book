.PHONY:mock
mock:
	@mockgen -source=internal/service/user.go -package=mock -destination=internal/service/mock/user.mock.go
	@mockgen -source=internal/service/code.go -package=mock -destination=internal/service/mock/code.mock.go
	@mockgen -source=internal/repository/user.go -package=mock -destination=internal/repository/mock/user.mock.go
	@mockgen -source=internal/repository/code.go -package=mock -destination=internal/repository/mock/code.mock.go
	@mockgen -source=internal/repository/dao/user.go -package=mock -destination=internal/repository/dao/mock/user.mock.go
	@mockgen -source=internal/repository/cache/user.go -package=mock -destination=internal/repository/cache/mock/user.mock.go
	@go mod tidy
.PHONY:wire
wire:
	@cd wire&&wire
