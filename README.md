# gRPC
## Protocol Buffer
### 1. Overview
- `Protocol Buffers` là một cơ chế trung lập về ngôn ngữ và nền tảng để tuần tự hóa (serializing) dữ liệu có cấu trúc.
### 2. Message
- `Message` là thành phần cơ bản được sử dụng để định nghĩa cấu trúc dữ liệu. Mỗi `message` là một cấu trúc dữ liệu có thể chứa nhiều trường (`fields`) với các kiểu dữ liệu khác nhau.
    ```proto
    syntax = "proto3";

    message SearchRequest {
        string query = 1;
        int32 page_number = 2;
        int32 results_per_page = 3;
    }
    ```

### 3. Generated code
- Để `generate code` từ file `.proto` ta cần tải `plugin`
    ```
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    ```
- Sau đó chạy câu lệnh sau để `compile`
    ```
    protoc --proto_path=IMPORT_PATH --go_out=DST_DIR path/to/file.proto
    ```

## gRPC
### 1. Overview
- `gRPC` là một framework `RPC` (`Remote Procedure Call` - Gọi Thủ tục Từ xa) hiện đại, mã nguồn mở, hiệu năng cao, có thể chạy trong mọi môi trường.
- Trong `gRPC`, `client` có thể gọi trực tiếp các phương thức trên một `server` đặt ở một máy khác như thể một đối tượng cục bộ, giúp ta dễ dàng tạo các ứng dụng và dịch vụ phân tán.
- Mặc định, `gRPC` sẽ sử dụng `protocol buffer` để định nghĩa cấu trúc dữ liệu và các phương thức dịch vụ `service`.
### 2. Service
- `gRPC` sử dụng `protocol buffer` để định nghĩa các `service` và các phương thức trong `service` đó.
    ```proto
    service HelloService {
        rpc SayHello (HelloRequest) returns (HelloResponse);
    }

    message HelloRequest {
        string greeting = 1;
    }

    message HelloResponse {
        string reply = 1;
    }
    ```
- `gRPC` cho phép chúng ta sử dụng 4 loại phương thức
    - `Unary RPC`: `client` sẽ gửi đúng 1 `request` đến `server` và nhận lại đúng 1 `response`, giống như 1 hàm bình thường
        ```proto
        rpc SayHello(HelloRequest) returns (HelloResponse);
        ```
    - `Server streaming RPC`: `client` gửi 1 `request` đến `server` và nhận 1 `stream` gồm nhiều `response` trả về từ `server`. `Client` sẽ đọc `message` từ `stream` này cho đến khi không còn `message` trả về nữa.
        ```proto
        rpc LotsOfReplies(HelloRequest) returns (stream HelloResponse);
        ```
    - `Client streaming RPC`: `client` sẽ gửi một loạt các `message` một cách tuần tự sang cho `server` sử dụng `stream` được cung cấp. `Server` sẽ trả về đúng một `response` sau khi đã đọc và xử lý hết các `request`.
        ```proto
        rpc LotsOfGreetings(stream HelloRequest) returns (HelloResponse);
        ```
    - `Bidirectional streaming RPC`: `client` và `server` sẽ sử dụng các `read-write stream` để gửi và nhận một dãy các `message`. 
        ```proto
        rpc BidiHello(stream HelloRequest) returns (stream HelloResponse);
        ```
### 3. gRPC với Go
- Để sử dụng `gRPC` với `Go`, ta cần tải các `plugin` sau
    ```
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
    ```
    export PATH="$PATH:$(go env GOPATH)/bin"
    ```
- Sử dụng câu lệnh sau để `compile` và sinh code Go từ file `.proto`
    ```
    protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    routeguide/route_guide.proto
    ```
