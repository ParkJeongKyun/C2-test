# C2-test

## [훈련 목적] C2 및 암호화 시뮬레이션 테스트 환경

## [Educational Purpose] C2 & Encryption Simulation Test Environment

> ⚠️ **절대 함부로 실행하지 마세요 (CRITICAL WARNING)**
>
> - **한국어:** 본 리포지토리는 승인된 보안 훈련 및 연구 목적으로만 사용되어야 합니다. 위험성을 충분히 숙지하지 않은 상태에서 파일이나 스크립트를 실행하지 마십시오. 허가되지 않은 환경에서의 실행 및 오용에 대한 모든 책임은 실행자에게 있습니다.
> - **English:** This repository is strictly for authorized security training and research purposes only. DO NOT execute any files or scripts unless you fully understand the risks. The user bears all responsibility for unauthorized execution or misuse.

## 흐름도

<img width="746" height="524" alt="image" src="https://github.com/user-attachments/assets/c31ed998-1783-47dc-877d-bd7a488af86b" />

## 디렉토리 구조

```bash
├── RANSOMNOTE_SAMPLE.txt           # 랜섬노트 샘플(실제로 사용하지는 않습니다. 샘플이에요.)
├── README.md                       # 설명(안내문서)
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
├── target                          # 샘플용 테스트 타켓 폴더(암호화 대상)
│   ├── sample.db                       # 샘플 DB
│   ├── sample.mov                      # 샘플 영상
│   ├── sample.ogg                      # 샘플 오디오
│   ├── sample.txt                      # 샘플 문서
│   └── sample.webp                     # 샘플 이미지
└── vba                             # VBA 스크립트 원본
    ├── attack.bas                      # hta에 삽입되어 있는 스크립트 원본
    └── macro.bas                       # 엑셀 파일내에 매크로로 삽입되어 있는 스크립트 원본
```

# 실행/빌드 플랫폼 환경

```bash
- Windows 64비트(x64) 환경
(매크로 및 exe 파일 실행 플랫폼)

# Go 파일을 직접 실행하는 경우에는 플랫폼은 상관없습니다!
```

# 조건

```bash
의도하지 않은 경로의 파일 손상 및 확산을 방지하기 위해서 다음과 같은 조건이 적용되어 있습니다.
1. 경로 제한
    - 실행 위치를 기준으로 1Depth 상위 디렉토리까지만 탐색을 수행합니다.
2. 대상 폴더 제한
    - 'target'이라는 이름의 폴더의 내부 파일만 대상으로 합니다.
3. 확장자 제한
    - 지정된 특정 확장자 목록에 대해서만 암호화를 수행합니다.
```

# 실행 방법

```bash
1. 엑셀 파일 매크로를 통해서 실행하는 방법
    1-1. 안티바이러스를 종료한다.(!!!책임 지지 않습니다!!!)
    1-2. 엑셀 파일을 좌클릭하여 매크로 실행을 허용한다.
    1-3. 엑셀 파일을 실행하여 매크로 실행을 허용한다.
2. hta를 mshta.exe로 실행하는 방법
    2-1. CMD에서 mshta.exe "%cd%\build\attack.hta"
    2-2. 파워셀에서 mshta.exe "$PWD\build\attack.hta"
3. 실행파일을 직접 실행하여 실행하는 방법
    3-1. build 디렉토리에 있는 encrypt.exe를 직접 실행한다.
```

# 엑셀 매크로 실행 방법
1. 안티 바이러스 사용시 자동 차단되기 때문에 비활성화 해주세요 (비활성화로 인한 피해는 책임지지 않습니다!, 사진은 HP wolf security)
<img width="150" height="100" alt="image" src="https://github.com/user-attachments/assets/0dff30ba-f60d-46a3-b39c-1d4ad014a0dd" />
<img width="150" height="100" alt="image" src="https://github.com/user-attachments/assets/73ad4a23-64df-4c21-8b33-be368839fbcb" />

3. 매크로가 포함된 엑셀을 실행하여 편집을 사용(순서중요!)
<img width="200" height="50" alt="image" src="https://github.com/user-attachments/assets/2704a547-7838-41fb-8344-900fd93595a2" />

4. 파일 속성에서 보안 차단 해제
<img width="100" height="150" alt="image" src="https://github.com/user-attachments/assets/6c949b39-b1cc-4b4a-baa0-761bacf33903" />
<img width="400" height="100" alt="image" src="https://github.com/user-attachments/assets/653f131a-d457-4acd-9c70-3f4d0bea905a" />

5. 엑셀에서 컨텐츠 매크로를 허용하여 직접실행
<img width="300" height="50" alt="image" src="https://github.com/user-attachments/assets/b5c5d9b0-6bc5-4f88-bcef-21d3faacf53f" />


## 아래와 같이 윈도우 디펜더에서 차단될수도 있습니다!
<img width="200" height="250" alt="image" src="https://github.com/user-attachments/assets/6cf12c84-149b-4131-960f-969e47c0f547" />
<img width="200" height="100" alt="image" src="https://github.com/user-attachments/assets/46fe4a4b-0d08-431f-ade2-c2a84cea1533" />
<img width="300" height="250" alt="image" src="https://github.com/user-attachments/assets/cbf00d9c-78c7-4980-9aae-25ee49896760" />


# Go 코드 직접 실행

```bash
# 암호화 테스트
go run ./go/encrypt.go
# 복호화 테스트
go run ./go/decryption.go
```

# EXE 빌드 방법

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

# 추가 개선 가능한 것들

```bash
## 허가되지 않은 환경에서 실행하시면 안됩니다.
- 암호화 시에 멀티 쓰레딩, 프로세싱을 통해서 암호화 하는 속도를 높일 수 있음
- 레파지토리에는 암호화 시켜서 업로드하여 악의적으로 리버싱을 하기 힘들도록 할 수 있음
- 디스크 드롭 없이 메모리상에서 직접 페이로드를 구동하게하는 방식을 적용할 수 있음
```
