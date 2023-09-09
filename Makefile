all:
	echo 'use compile task'

deploy:
	ssh krownguards@136.243.176.153 "source ~/.profile ;\
		cd /var/www/vhosts/krownguards.com/go/src/ws/ ;\
		git pull --rebase ;\
		go get ;\
		go build"
