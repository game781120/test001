docker build -t digital-visitor-images .
docker run -itd --network host --privileged=true --restart=always --name digital-visitor-container digital-visitor-images:latest