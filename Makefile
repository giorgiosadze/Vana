build:
	podman build -t vana -f DockerFile .
start:
	podman run -d --rm -p 8000:8000 vana:latest
stop:
	podman stop $(shell podman ps -q --filter ancestor=vana)
clean:
	podman rmi vana

.PHONY: build run clean
