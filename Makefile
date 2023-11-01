# Define the name of the executable binary.
BINARY_NAME := prostoprogram

all: build gooseup run clean

build:
	go build -o $(BINARY_NAME) main.go

run:
	./$(BINARY_NAME)

gooseup:
	goose -dir migrations mysql "kursUser:kursPswd@tcp(127.0.0.1:3306)/TEST" up

goosedown:
	goose -dir migrations mysql "kursUser:kursPswd@tcp(127.0.0.1:3306)/TEST" down

clean:
	if [ -f $(BINARY_NAME) ] ; then rm $(BINARY_NAME) ; fi

.PHONY:
	all build run clean
