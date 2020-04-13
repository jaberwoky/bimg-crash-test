test:
	docker build -t bimg-crash-test .
	docker run -it --rm bimg-crash-test
