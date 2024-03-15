build:
	docker build -t ttl.sh/bd/gohub/v1:1h .

run:
	docker run -it --rm --name gohub \
	-d \
	-p 8000:8000 \
	-p 9000:9000 \
 	--env branch=owen/swagger \
    --env repository=https://xxx:xxxx@github.com/bytedance-soft/wallet-server.git \
    --env commonRepository=https://xxx:xxxx@github.com/bytedance-soft/go-common.git \
 	ttl.sh/bd/gohub/v1:1h

delete:
	docker stop gohub && docker rm gohub
