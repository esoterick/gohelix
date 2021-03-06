package gohelix

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"code.google.com/p/go.crypto/ssh"
)

// AddTestCluster calls helix-admin.sh --zkSvr localhost:2181 --addCluster
func AddTestCluster(cluster string) error {
	cmd := "/opt/helix/bin/helix-admin.sh --zkSvr localhost:2181 --addCluster " + strings.TrimSpace(cluster)
	if _, err := RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

// AddNode /opt/helix/bin/helix-admin.sh --zkSvr localhost:2181  --addNode
func AddNode(cluster string, host string, port string) error {

	cmd := fmt.Sprintf("/opt/helix/bin/helix-admin.sh --zkSvr localhost:2181  --addNode %s %s:%s", cluster, host, port)
	if _, err := RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

// AddResource /opt/helix/bin/helix-admin.sh --zkSvr localhost:2181 --addResource
func AddResource(cluster string, resource string, replica string) error {
	cmd := fmt.Sprintf("/opt/helix/bin/helix-admin.sh --zkSvr localhost:2181 --addResource %s %s %s MasterSlave", cluster, resource, replica)
	if _, err := RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

// Rebalance /opt/helix/bin/helix-admin.sh --zkSvr localhost:2181 --rebalance
func Rebalance(cluster string, resource string, replica string) error {
	cmd := fmt.Sprintf("/opt/helix/bin/helix-admin.sh --zkSvr localhost:2181 --rebalance %s %s %s", cluster, resource, replica)
	if _, err := RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

// DropTestCluster /opt/helix/bin/helix-admin.sh --zkSvr localhost:2181 --dropCluster
func DropTestCluster(cluster string) error {
	cmd := "/opt/helix/bin/helix-admin.sh --zkSvr localhost:2181 --dropCluster " + strings.TrimSpace(cluster)
	if _, err := RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

// StartController sudo /usr/bin/supervisorctl start helixcontroller
func StartController() error {
	if _, err := RunCommand("sudo /usr/bin/supervisorctl start helixcontroller"); err != nil {
		return err
	}
	return nil
}

// StopController sudo /usr/bin/supervisorctl stop helixcontroller
func StopController() error {
	if _, err := RunCommand("sudo /usr/bin/supervisorctl stop helixcontroller"); err != nil {
		return err
	}
	return nil
}

// StartParticipant /usr/bin/supervisorctl start participant_
func StartParticipant(port string) error {
	command := "/usr/bin/supervisorctl start participant_" + port
	if _, err := RunCommand(command); err != nil {
		return err
	}

	return nil
}

// StopParticipant /usr/bin/supervisorctl stop participant_
func StopParticipant(port string) error {
	command := "/usr/bin/supervisorctl stop participant_" + port
	if _, err := RunCommand(command); err != nil {
		return err
	}

	return nil
}

// RunCommand execute command via ssh
func RunCommand(command string) (string, error) {
	key, err := getKeyFile()
	if err != nil {
		return "", err
	}

	config := &ssh.ClientConfig{
		User: "vagrant",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	client, err := ssh.Dial("tcp", "127.0.0.1:2222", config)

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run(command); err != nil {
		return "", err
	}
	return b.String(), nil
}

func getKeyFile() (key ssh.Signer, err error) {
	out, err := exec.Command("/usr/bin/vagrant", "ssh-config").Output()
	if err != nil {
		return
	}

	identityFile := ""
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, "IdentityFile") {
			parts := strings.Fields(line)
			identityFile = parts[1]
			break
		}
	}

	if identityFile == "" {
		return
	}

	buf, err := ioutil.ReadFile(identityFile)
	if err != nil {
		return
	}
	key, err = ssh.ParsePrivateKey(buf)
	if err != nil {
		return
	}
	return
}
