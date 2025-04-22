APPS_PATH = apps
GENERATED_PATH = proto

PROTOC = protoc

GO_INSTALL_COMMANDS = \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

MICROSERVICES = auth

# Правило для генерации всех файлов
all: $(MICROSERVICES)

# Правила для каждого микросервиса
$(MICROSERVICES):
	@echo "Building protobuf files for $@..."
	@mkdir -p $(APPS_PATH)/$@/$(GENERATED_PATH)
	$(PROTOC) --proto_path=$(APPS_PATH)/$@/proto \
	           --go_out=$(APPS_PATH)/$@/$(GENERATED_PATH) \
	           --go-grpc_out=$(APPS_PATH)/$@/$(GENERATED_PATH) \
	           $(APPS_PATH)/$@/proto/*.proto
	@echo "Protobuf files generated successfully for $@."

# Установка плагинов для protoc
install:
	$(GO_INSTALL_COMMANDS)

# Очистка сгенерированных файлов
clean:
	find $(APPS_PATH) -name '*.pb.go' -delete

.PHONY: all install clean $(MICROSERVICES)
