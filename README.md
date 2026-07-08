# ⚠️ 절때로 함부로 실행하지마세요 ⚠️

# ⚠️ WARNING: DO NOT execute any files unless you fully understand the risks. ⚠️

# C2-test

본 레파지토리는 훈련용 랜섬웨어 레파지토리입니다.
실제 소스 코드 형상관리겸 C2 페이로드 드롭퍼 및 키 수신 채널로 사용됩니다.

This repository is strictly for educational and testing purposes regarding ransomware simulations (C2 Testing).

<img width="1458" height="982" alt="image" src="https://github.com/user-attachments/assets/4cf202e7-3b92-4d77-ab06-6d5eb7c3fe8d" />

# 디렉토리 구조

```bash
├── RANSOMNOTE_SAMPLE.txt           # 랜섬노트 샘플(실제로 사용하지 않습니다)
├── README.md                       # 설명
├── build                           # 빌드된 파일
│   ├── attack.hta                      # mshta.exe로 실행할 마크다운 파일(VBA가 삽입되어 있음)
│   ├── decryption.exe                  # 복호화용 프로그램(윈도우용)
│   ├── encoded_attack.txt              # attack.hta를 Base64로 인코딩한 파일
│   ├── encoded_encrypt.txt             # encrypt.exe를 Base64로 인코딩한 파일
│   └── encrypt.exe                     # [실제랜섬웨어기능]암호화용 프로그램(윈도우용)
├── go                              # 원본 소스코드
│   ├── decryption.go                   # 복호화용 원본 소스코드
│   ├── encrypt.go                      # 암호화용 원본 소스코드
│   └── go.mod
├── key                             # 암/복호화에 사용할 RSA키
│   ├── private.pem                     # 개인키(복호화용)
│   └── public.pem                      # 공개키(암호화용)
├── macro                           # VBA가 삽입된 엑셀파일
│   └── [최종][실행X](가상기업)_월차_신청서_양식.xlsm
└── target                          # 샘플용 테스트 타켓 폴더(암호화 대상)
    ├── sample.db                       # 샘플 DB
    ├── sample.mov                      # 샘플 영상
    ├── sample.ogg                      # 샘플 오디오
    ├── sample.txt                      # 샘플 문서
    └── sample.webp                     # 샘플 이미지
```

# 조건

```bash
안전을 위해서 다음과 같은 조건이 적용되어 있습니다.
1. 1Depth 상위 디렉토리까지만 검사함.
2. 'target'이라는 이름의 폴더만 선택함.
3. 그 폴더 내 특정 확장자 별로만 암호화 수행함.
```

# 실행방법

```bash
1. 엑셀 파일 매크로를 통해서 실행하는 방법
    1-1. 안티바이러스를 종료한다.(책임은 지지 않습니다!!!)
    1-2. 엑셀 파일을 좌클릭하여 매크로 실행을 허용한다.
    1-3. 엑셀 파일을 실행하여 매크로 실행을 허용한다.
2. hta를 mshta.exe로 실행하는 방법
    2-1. CMD에서 mshta.exe "%cd%\build\attack.hta"
    2-2. 파워셀에서 mshta.exe "$PWD\build\attack.hta"
3. 실행파일을 직접 실행하여 실행하는 방법
```

#

# 테스트 고 코드 실행

```bash
# 암호화 테스트
go run ./go/encrypt.go
# 복호화 테스트
go run ./go/decryption.go
```

# 빌드 방법

```bash
# 실제로 빌드시에는 난독화가 필요하겠지만 훈련용이기 때문에 사용하지 않습니다.
# 윈도우 기준 빌드
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./build/encrypt.exe ./go/encrypt.go
# 복호화 윈도우 기준 빌드
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./build/decryption.exe ./go/decryption.go
# base64 인코딩
base64 -i encrypt.exe -o encoded_encrypt.txt
base64 -i attack.hta -o encoded_attack.txt
```
