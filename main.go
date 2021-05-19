package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
	"k8s.io/kubernetes/test/utils"
)

var RootCmd = &cobra.Command{
	Use:   "self-signed-cert",
	Short: "Create a self-signed TLS certificate.",
	Args:  cobra.MaximumNArgs(0),
	Run:   run,
}

var (
	namespace string
	service   string
)

func init() {
	RootCmd.Flags().StringVar(&namespace, "namespace", "",
		"Namespace in which the service resides into.")
	RootCmd.Flags().StringVar(&service, "service-name", "",
		"Service for which to generate the certificate.")
	RootCmd.MarkFlagRequired("namespace")
	RootCmd.MarkFlagRequired("service-name")
}

// Setup the server cert. For example, user apiservers and admission webhooks
// can use the cert to prove their identify to the kube-apiserver
func setupServerCert(namespaceName, serviceName string) {
	certDir, err := ioutil.TempDir("", "self-signed-certificate")
	if err != nil {
		log.Fatalf("Failed to create a temp dir for cert generation %v", err)
	}

	signingKey, err := utils.NewPrivateKey()
	if err != nil {
		log.Fatalf("Failed to create CA private key %v", err)
	}
	signingCert, err := cert.NewSelfSignedCACert(cert.Config{CommonName: "e2e-server-cert-ca"}, signingKey)
	if err != nil {
		log.Fatalf("Failed to create CA cert for apiserver %v", err)
	}

	caCertFile := filepath.Join(certDir, "ca.crt")

	if err := ioutil.WriteFile(caCertFile, utils.EncodeCertPEM(signingCert), 0644); err != nil {
		log.Fatalf("Failed to write CA cert %v", err)
	}
	key, err := utils.NewPrivateKey()
	if err != nil {
		log.Fatalf("Failed to create private key for %v", err)
	}
	signedCert, err := utils.NewSignedCert(
		&cert.Config{
			CommonName: serviceName + "." + namespaceName + ".svc",
			AltNames:   cert.AltNames{DNSNames: []string{serviceName + "." + namespaceName + ".svc"}},
			Usages:     []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		key, signingCert, signingKey,
	)
	if err != nil {
		log.Fatalf("Failed to create cert%v", err)
	}

	certFile := filepath.Join(certDir, "server.crt")
	keyFile := filepath.Join(certDir, "server.key")

	if err = ioutil.WriteFile(certFile, utils.EncodeCertPEM(signedCert), 0600); err != nil {
		log.Fatalf("Failed to write cert file %v", err)
	}

	privateKeyPEM, err := keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		log.Fatalf("Failed to marshal key %v", err)
	}

	if err = ioutil.WriteFile(keyFile, privateKeyPEM, 0644); err != nil {
		log.Fatalf("Failed to write key file %v", err)
	}

	fmt.Printf("%s\n", certDir)
}

func run(cmd *cobra.Command, args []string) {
	setupServerCert(namespace, service)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
