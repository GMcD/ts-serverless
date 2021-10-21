# This makes the subsequent variables available to child shells
.EXPORT_ALL_VARIABLES:

include .env

# Collect Last Target, convert to variable, and consume the target.
# Allows passing arguments to the target recipes from the make command line.
CMD_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# Consume them to prevent interpretation as targets
$(eval $(CMD_ARGS):;@:)
# Service for command args
ARGUMENT  := $(word 1,${CMD_ARGS})

##
## Usage:
##  make [target] [ARGUMENT]
##   operates in namespace ${NAMESPACE}
##

help:		## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

stack:		## Update ECR tags in stack.yml on main
	git checkout main && \
	if (( $$(git status --porcelain | wc -l) > 0 )); then \
	    printf "$${GREEN}Module $${RED}ts-serverless$${GREEN} has changes, run $${CYAN}make commit <message>$${GREEN} first.$${NC}\n"; \
	    exit 1; \
	fi && \
	awk -F "." '/354455067292/ { printf $$1; for(i=2;i<NF;i++) printf FS$$i; print FS$$NF+1 } !/354455067292/ { print }' stack.yml > .stack.yml && mv .stack.yml stack.yml

bump:		## Update go mod version numbers on main
bump: stack
	npm --no-git-tag-version version patch && \
	for mod in $$(find ./micros -name \*.mod); do \
		awk -F "1." '/ts-serverless v/ { printf $$1; for(i=2;i<NF;i++) printf FS$$i; print FS$$NF+1 } !/ts-serverless v/ { print }' $$mod > $${mod}.tmp && mv $${mod}.tmp $$mod; \
	done  && \
	for micro in $$(ls -d micros/*/); do pushd ./$${micro}; go mod tidy; popd; done && \
	git add . ; git commit -m Version-$$(cat package.json | jq -j '.version'); git push

commit:		## Short hand for Commit to prod
	git add .; git commit -m ${ARGUMENT}; git push

fork:		## Short hand for Commit main to Fork Remote
fork: bump
	git checkout gmcd && \
	git merge main && \
	git add . ; git commit -m ${ARGUMENT}; git push fork HEAD:master 

tag:		## Tag a Release
tag: fork 
	git tag v$$(cat package.json | jq -j '.version') -am ${ARGUMENT} && \
	git push fork HEAD:master --tags && \
	git checkout main

logs:		## Log Pod ${ARGUMENT} by prefix
logs:
	kubectl logs --namespace openfaas-fn $(shell kubectl get pods --namespace openfaas-fn -o=jsonpath='{.items[*].metadata.name}' -l faas_function=${ARGUMENT})

login:  	## ECR Docker Login
	@ aws ecr get-login-password --region $${AWS_REGION} | docker login --username AWS --password-stdin $${AWS_ACCOUNT_ID}.dkr.ecr.$${AWS_REGION}.amazonaws.com
	@ AWS_ACCOUNT_ID=$$(aws sts get-caller-identity --output text --query 'Account'); \
	AWS_IAM_ARN=$$(aws sts get-caller-identity --output text --query 'Arn'); \
	echo "Running as $${AWS_IAM_ARN} in $${AWS_REGION} for $${AWS_ACCOUNT_ID}."

up:		## Run FaaS up
up:			# If a build breaks due to syntax errors, there is an issue with pseudo-version package
up:			# numbers, whereby the error is uploaded to public github, and subsequent builds will
up:			# pick up the new release, when it is pushed, as it can't be pushed til it is built..
up:			# Current resolution is to rebuild the broken guy, out of band, with
up:			# `docker builder prune -f` && \
up:			# `GOPRIVATE=github.com/GMcD faas build --no-cache --build-arg GO111MODULE=on --filter collective` 
up: login
	for micro in $$(ls -d micros/*/); do pushd ./$${micro}; GOPRIVATE=github.com/GMcD go mod tidy; popd; done
	# Update micros with new core && code bases
	# ./update-micros.sh telar-core
	# ./update-micros.sh telar-web
	echo "Running FaaS up..."
	GOPRIVATE=github.com/GMcD faas up --build-arg GO111MODULE=on
