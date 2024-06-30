package utils

import (
	"os"
	"io"
	"fmt"
	"time"
	"context"
	"golang.org/x/crypto/ssh"
    "github.com/bramvdbogaerde/go-scp"

    "root/config"
)


func SSHConnect(user, password, host string, port int) (*ssh.Client, error) {
    config := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            ssh.Password(password),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         30 * time.Second,
    }

    address := fmt.Sprintf("%s:%d", host, port)
    return ssh.Dial("tcp", address, config)
}

func SSHRunCommand(client *ssh.Client, command string) (string, error) {
    session, err := client.NewSession()
    if err != nil {
        return "", fmt.Errorf("failed to create session: %w", err)
    }
    defer session.Close()

    var stdout, stderr io.Reader
    stdout, err = session.StdoutPipe()
    if err != nil {
        return "", fmt.Errorf("failed to capture stdout: %w", err)
    }
    stderr, err = session.StderrPipe()
    if err != nil {
        return "", fmt.Errorf("failed to capture stderr: %w", err)
    }

    if err := session.Start(command); err != nil {
        return "", fmt.Errorf("failed to start command: %w", err)
    }

    outBytes, err := io.ReadAll(stdout)
    if err != nil {
        return "", fmt.Errorf("failed to read stdout: %w", err)
    }
    errBytes, err := io.ReadAll(stderr)
    if err != nil {
        return "", fmt.Errorf("failed to read stderr: %w", err)
    }

    if err := session.Wait(); err != nil {
        return "", fmt.Errorf("command failed: %w\nstderr: %s", err, string(errBytes))
    }

    return string(outBytes), nil
}

func FileTransfer(sshClient *ssh.Client, localFilePath string, remoteFilePath string) error {
    client, _ := scp.NewClientBySSH(sshClient)
    defer client.Close()

    f, _ := os.Open(localFilePath)
	defer f.Close()

    err := client.CopyFile(context.Background(), f, remoteFilePath, "0666")
	if err != nil {
		return err
	}

    return nil
}

func GetTables(databaseName string) ([]config.Table){
	return config.DatabaseTableMap[databaseName]
}

func Contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}