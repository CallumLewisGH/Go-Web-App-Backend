package integration_test_config

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"testing"

	"github.com/joho/godotenv"
)

type DockerContainer struct {
	ID        string
	Name      string
	Port      int
	Cleanup   func()
	DbConnStr string
}

func StartTestDatabaseContainer(t *testing.T) *DockerContainer {
	t.Helper()

	port, err := GetAvailablePort()
	if err != nil {
		t.Fatalf("Failed to get available port: %v", err)
	}

	err = godotenv.Load("/home/callum/Desktop/Go-Web-App-Backend/.test.env")
	if err != nil {
		t.Fatalf("Failed to Load Env: %v", err)
	}

	containerName := fmt.Sprintf("%s-%d", os.Getenv("DOCKER_CONTAINER_NAME"), port)
	args := []string{
		"run", "-d",
		"--name", containerName,
		"-e", fmt.Sprintf("POSTGRES_USER=%s", os.Getenv("DATABASE_USER")),
		"-e", fmt.Sprintf(`POSTGRES_PASSWORD=%s`, os.Getenv("DATABASE_PASSWORD")),
		"-e", fmt.Sprintf("POSTGRES_DB=%s_%d", os.Getenv("DATABASE_NAME"), port),
		"-p", fmt.Sprintf("%d:5432", port),
		"postgres:latest",
	}

	cmd := exec.Command("docker", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	containerID := string(bytes.TrimSpace(out.Bytes()))

	connStr := fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s_%d sslmode=disable",
		port,
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		port,
	)

	return &DockerContainer{
		ID:   containerID,
		Port: port,
		Cleanup: func() {
			exec.Command("docker", "rm", "-f", "-v", containerID).Run()
		},
		DbConnStr: connStr,
	}
}
func GetAvailablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}
