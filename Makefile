deploy-namespace:
	@kubectl apply -f deploy/kubernetes/simple-ms-namespace.yaml --wait=true

deploy-all: deploy-namespace
	@for file in $(shell ls -I simple-ms-namespace.yaml deploy/kubernetes) ; do \
        kubectl apply -f deploy/kubernetes/$$file ; \
    done

remove-deployment:
	@kubectl delete all,ing -n nsch-ms --all && kubectl delete ns nsch-ms