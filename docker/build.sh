tar cvf hola2-mv-consumer.tar hola2-mv-consumer
docker build -t registry.tde.sktelecom.com/sktsdn/sdn-bp-dockerfiles/hola2-mv-consumer:v2.240820.170525 -f docker/Dockerfile . --platform linux/amd64
