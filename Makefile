build:
	$(MAKE) -C user-service
	$(MAKE) -C user-cli
	$(MAKE) -C vessel-service
	$(MAKE) -C consignment-service
	# $(MAKE) -C consignment-cli
