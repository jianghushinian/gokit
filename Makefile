.PHONY: vet
vet:
	@echo "---- Vetting ----"
	go vet ./...
	@echo "---- Successfully Vet ----\n"


.PHONY: test
test:
	@echo "---- Testing ----"
	go test -count=1 -v -cover -p 1 ./...
	@echo "---- Successfully Tested ----\n"
