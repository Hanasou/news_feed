# Build all images

# Gateway
docker build -f go/gateway/Dockerfile -t news_feed_gateway:latest .

# Users service
docker build -f go/user/Dockerfile -t news_feed_user_service:latest .