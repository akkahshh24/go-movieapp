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

testrating1:
	grpcurl -plaintext -d '{"record_id":"1", "record_type":"movie"}' localhost:8082 RatingService/GetAggregatedRating

movie1:
	cd movie/cmd && go run main.go --port=8083

movie2:
	cd movie/cmd && go run main.go --port=8086

movie3:
	cd movie/cmd && go run main.go --port=8089

consul:
	docker run -d -p 8500:8500 -p 8600:8600/udp --name dev-consul hashicorp/consul agent -server -ui -node=server-1 -bootstrap-expect=1 -client='0.0.0.0'

proto:
	protoc -I=api --go_out=. --go-grpc_out=. movie.proto

benchmark:
	cd cmd/sizecompare && go test -bench=.

.PHONY: metadata1 metadata2 metadata3 rating1 rating2 rating3 movie1 movie2 movie3 testrating1 consul proto benchmark