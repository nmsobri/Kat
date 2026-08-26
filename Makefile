.PHONY: clean

kat:
	@echo "Building Kat..."
	@go build -o kat

clean:
	@echo "Cleaning..."
	@rm -rf kat
