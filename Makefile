NAME       = series-dl
REPO_URL   = github.com/moolen/$(NAME)

GO         = go
BUILDFLAGS = -v -installsuffix cgo -o $(NAME) --tags release

M          = $(shell printf ">>>")

.PHONY: all vendor clean
all: $(GITSEMVER) vendor
	@echo -n $(info $(M) building binary)
	CGO_ENABLED=0 $(GO) build $(BUILDFLAGS)

vendor: 
	@echo -n $(info $(M) installing deps)
	dep ensure -v
	@touch $@

clean:
	@echo -n $(info $(M) cleaning up)
	@rm -vf $(NAME)
