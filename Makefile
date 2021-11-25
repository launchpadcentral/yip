run: 
	# for developing
	docker-compose up yip

test:  
	# test before building
	docker-compose up yip-test

build:	
	# set GOOS, GOARCH and RELEASE first then build
	docker-compose up yip-build

exec:
	# for developing
	docker-compose exec yip /bin/bash

clean:	
	# quick cleanup
	docker-compose down
