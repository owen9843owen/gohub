build:
	docker build -t ttl.sh/bd/gohub/v1:1h .

run:
	docker run -it --rm --name gohub \
	-d \
 	--env branch=owen/swagger \
    --env repository=https://owen9843owen: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@github.com/bytedance-soft/wallet-server.git \
    --env server=wallet-server \
    --env commonRepository=https://owen9843owen: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@github.com/bytedance-soft/go-common.git \
 	ttl.sh/bd/gohub/v1:1h

delete:
	docker stop gohub && docker rm gohub
