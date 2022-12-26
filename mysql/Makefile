## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

build: ## - build the source (latest)
	@docker compose build --build-arg GIT_COMMIT=`git rev-parse HEAD` --build-arg GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
	@docker image prune -f

run: ## - run bludgeon-mysql
	@docker container rm -f mysql
	@docker image prune -f
	@docker compose up -d

stop: ## - stop bludgeon-mysql
	@docker compose down

clean: stop ## - stop and clean docker resources
	@docker compose rm -f