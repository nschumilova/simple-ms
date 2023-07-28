define deploy_postgres
	helm upgrade postgresql bitnami/postgresql --install --values deploy/helm/postgres/values.yaml \
	 --namespace nsch-db --create-namespace --wait
endef

define deploy_simplems
	helm upgrade simplems-users deploy/helm/simple-ms/ --install --namespace nsch-ms --create-namespace --wait
endef

define uninstall_postgres
	helm delete postgresql --namespace nsch-db --wait
	kubectl delete pvc -l app.kubernetes.io/name=postgresql -n nsch-db --wait --timeout=60s
	kubectl delete ns nsch-db --wait --timeout=60s
endef

define uninstall_simplems
	helm delete simplems-users --namespace nsch-ms --wait
	kubectl delete ns nsch-ms --wait --timeout=60s
endef



build-pg-schema:
	DOCKER_BUILDKIT=1 docker image build -t nshumilova/simple-ms-pgdb-schema:$(version) -f build/pgschema.Dockerfile .

push-pg-schema:
	docker image push nshumilova/simple-ms-pgdb-schema:$(version)

build-serice:
	DOCKER_BUILDKIT=1 docker image build -t nshumilova/simple-ms:$(version) -f build/simplems.Dockerfile .

push-service:
	docker image push nshumilova/simple-ms:$(version)



add-postgres-repository:
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update bitnami

deploy-postgres: add-postgres-repository
	$(call deploy_postgres)

deploy-service: 
	$(call deploy_simplems)

deploy-all: 
	$(call deploy_postgres)
	$(call deploy_simplems)

uninstall-postgres:
	$(call uninstall_postgres)

uninstall-service:
	$(call uninstall_simplems)

uninstall-all:
	$(call uninstall_simplems)
	$(call uninstall_postgres)