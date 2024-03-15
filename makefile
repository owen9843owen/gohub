build:
	docker build -t ttl.sh/bd/gohub/v1:1h .

run:
	docker run -it --rm --name gohub \
	-d \
 	--env branch=owen/swagger \
    --env path=/Users/owen/Documents/project/github.com/bytedance-soft/wallet-server/ \
    --env port=8088 \
    --env repository=https://owen9843owen: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@github.com/bytedance-soft/wallet-server.git \
    --env server=wallet-server \
    --env commonRepository=https://owen9843owen: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@github.com/bytedance-soft/go-common.git \
 	ttl.sh/bd/gohub/v1:1h
