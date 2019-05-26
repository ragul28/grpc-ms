all:
	(cd user-service; make build)
	(cd user-cli; make build)
	(cd vessel-service; make build)
	(cd consignment-service; make build)
	# (cd consignment-cli; make build)


