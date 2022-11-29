MOCKS_DESTINATION=mocks
FILES = internal/transport/rest/server.go \
		internal/transport/rest/handler/shortener.go \
		internal/transport/rest/handler/handler.go \
		internal/storage/sqllite/storage.go \
		internal/service/shortener.go \
		internal/core/model.go


create: 
	@echo "Generating mocks...!!!"
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $(FILES); do \
	echo $$file; \
	mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file;  \
	done


# mocks: 
# 	@echo "Generating mocks...!!!"
# 	@rm -rf $(MOCKS_DESTINATION)

# 	@for file in 1 2 3 4 ; do  echo file; 
#do mockgen -source=$$file -destination=$(MOCKS_DESTINATION) 



# MOCKS_DESTINATION=mocks
# .PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed