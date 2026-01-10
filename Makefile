.PHONY: run bin clean

bin = kat

run: bin
	@./$(bin) -run doc/type.kat

bin:
	@go build .

clean:
	@rm -rf $(bin)