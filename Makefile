metadata1:
    cd metadata/cmd && go run main.go --port=8081

metadata2:
	cd metadata/cmd && go run main.go --port=8084

metadata3:
	cd metadata/cmd && go run main.go --port=8087

rating1:
	cd rating/cmd && go run main.go --port=8082

rating2:
	cd rating/cmd && go run main.go --port=8085

rating3:
	cd rating/cmd && go run main.go --port=8088

testgetrating1:
	grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie"}' localhost:8082 RatingService/GetAggregatedRating

testputrating1:
	grpcurl -plaintext -d '{"record_id":"1", "record_type": "movie", "user_id": "alex", "rating_value": 5}' localhost:8082 RatingService/PutRating

movie1:
	cd movie/cmd && go run main.go --port=8083

movie2:
	cd movie/cmd && go run main.go --port=8086

movie3:
	cd movie/cmd && go run main.go --port=8089

consul:
	docker run -d -p 8500:8500 -p 8600:8600/udp --name dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client='0.0.0.0'

kafka:
	cd cmd/ratingproducer && docker compose up -d

create-topic:
	docker exec -it kafka kafka-topics.sh --zookeeper zookeeper:2181 --replication-factor 1 --partitions 1 --create --topic test-topic-1

producer:
	cd cmd/ratingproducer && go run main.go

mysql:
	docker run --name movieapp_db -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=movieapp -p 3306:3306 -d mysql:latest

create-tables:
	docker exec -i movieapp_db mysql -uroot -ppassword -D movieapp < schema/schema.sql

exec-mysql:
	docker exec -it movieapp_db mysql -uroot -ppassword -D movieapp

show-tables:
	docker exec -it movieapp_db mysql -uroot -ppassword -D movieapp -e "SHOW tables"

proto:
	protoc -I=api --go_out=. --go-grpc_out=. movie.proto

benchmark:
	cd cmd/sizecompare && go test -bench=.

mock:
	mockgen -source=metadata/internal/controller/metadata/controller.go -destination=gen/mock/metadata/repository/repository.go -package=repository

unit-test:
	go test -cover ./...

integration-test:
	go run test/integration/*.go

.PHONY: \
	metadata1 metadata2 metadata3 \
	rating1 rating2 rating3 \
	movie1 movie2 movie3 \
	testgetrating1 testputrating1 \
	consul kafka create-topic producer mysql create-tables exec-mysql show-tables \
	proto benchmark mock unit-test integration-test
