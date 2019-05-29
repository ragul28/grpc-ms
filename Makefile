build:
	$(MAKE) -C user-service
	$(MAKE) -C vessel-service
	$(MAKE) -C consignment-service
	$(MAKE) -C client-cli
