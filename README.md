# C2-test

랜섬웨어 실습을 위한 테스트 레파지토리입니다.
C2 테스트 입니다!
!!!! 절때로 함부로 실행하지마세요 !!!!

This repository is strictly for educational and testing purposes regarding ransomware simulations (C2 Testing).
⚠️ WARNING: DO NOT execute any files unless you fully understand the risks.

# 디렉토리 구조

```bash
├── README.md           - 설명 파일
├── build               - 빌드 결과 디렉토리(실행하시면 안됩니다! 테스트용 악성코드에요)
│   └── core.exe
├── go                  - 원본 파일
│   ├── decryption.go
│   ├── encrypt.go
│   └── go.mod
├── key                 - 암호화 키 파일
│   ├── private.pem     - 개인키(실제로 해커가 이걸 올려둘리는 없죠)
│   └── public.pem      - 공개키(암호화용 키에요)
└── target              - 암호화 대상
    ├── sample.db       - 샘플 디비파일
    ├── sample.mov      - 샘플 영상
    ├── sample.ogg      - 샘플 오디오
    ├── sample.txt      - 샘플 문서
    └── sample.webp     - 샘플 이미지
```

# 빌드 방법

```bash
cd ./go
go mod init c2-test-go
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o core.exe main.go
```
