# minder
Structure: Scrim Architecture Combined With Repository Pattern
In this project, have many layers (handler/controller in this case minder_handler.go > usecase for business logic in this case minder_usecase.go > repository for db access, and model as entity in this case minder.go)

Scrim is good for microservice application and repository pattern is good for implementing abstraction repository which is good for small application

Linting: i am using sonarlint (but golang has golint as built-in tool)

How to run the service:
1. make properties file (minder.properties, minder.yaml, minder.env, minder.ini, etc) in /opt/secret/ (for windows you can make in C drive) and run go mod tidy for getting all required third party libraries
2. change database config in properties file according to your specific database configuration
3. create table with this query
CREATE TABLE Users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number varchar(20) NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULl,
    gender ENUM('male', 'female') NOT NULL,
    date_of_birth DATE NOT NULL,
    profile_picture TEXT not NULL,
    is_upgraded BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE Swipes (
    swipe_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    target_user_id INT,
    swipe_action ENUM('like', 'pass') NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES Users(user_id),
    FOREIGN KEY (target_user_id) REFERENCES Users(user_id)
);
4. you can run
for development
"go run main.go" in minder project
for production:
Build:
-windows : GOOS=windows GOARCH=amd64 go build -o myapp.exe
-macos : GOOS=darwin GOARCH=amd64 go build -o myapp_darwin
-linux : GOOS=linux GOARCH=amd64 go build -o myapp_linux
-linux with ARM: GOOS=linux GOARCH=arm go build -o myapp_arm
Run:
-windows : run myapp.exe in command line
-macos : run ./myapp_darwin
-linux : run ./myapp_linux

list of goos and goarch (go tool dist list)
aix/ppc64
android/386
android/amd64
android/arm
android/arm64
darwin/amd64
darwin/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
freebsd/arm64
freebsd/riscv64
illumos/amd64
ios/amd64
ios/arm64
js/wasm
linux/386
linux/amd64
linux/arm
linux/arm64
linux/loong64
linux/mips
linux/mips64
linux/mips64le
linux/mipsle
linux/ppc64
linux/ppc64le
linux/riscv64
linux/s390x
netbsd/386
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/386
openbsd/amd64
openbsd/arm
openbsd/arm64
openbsd/ppc64
plan9/386
plan9/amd64
plan9/arm
solaris/amd64
wasip1/wasm
windows/386
windows/amd64
windows/arm
windows/arm64

5. untuk register gambar yang dimasukan sudah dalam bentuk base64 (https://emn178.github.io/online-tools/base64_encode_file.html)