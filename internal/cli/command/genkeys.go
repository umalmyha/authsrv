package command

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/umalmyha/authsrv/internal/cli/args"
)

type genKeysCommandOptions struct {
	help        bool
	privateFile string
	publicFile  string
}

type genKeysCommand struct {
	*LoggingCommand
	args args.ParsedArgs
}

func NewGenKeysCommand(args args.ParsedArgs, logger *log.Logger) Executor {
	return &genKeysCommand{
		LoggingCommand: &LoggingCommand{logger: logger},
		args:           args,
	}
}

func (c *genKeysCommand) Run() error {
	options := c.extractOptions()
	if options.help {
		c.Help()
		return nil
	}

	if options.privateFile == "" {
		options.privateFile = "private.pem"
	}

	if options.publicFile == "" {
		options.publicFile = "public.pem"
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateFile, err := os.Create(options.privateFile)
	if err != nil {
		return err
	}
	defer privateFile.Close()

	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	publicFile, err := os.Create(options.publicFile)
	if err != nil {
		return err
	}
	defer publicFile.Close()

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return err
	}

	logger := c.Logger()
	logger.Println("private and public key files are generated successfully")

	return nil
}

func (c *genKeysCommand) Help() {
	logger := c.Logger()
	logger.Println("genkeys - command create private and public key files for JWT generation")
	logger.Println("options:")
	logger.Println("  --help - show help")
	logger.Println("  --privateFile - specify filename for private key (default is private.pem)")
	logger.Println("  --publicFile - specify filename for public key (default is public.pem)")
	logger.Println("example:")
	logger.Println("  genkeys --privateFile=priv.pem --publicFile=pub.pem")
}

func (c *genKeysCommand) extractOptions() genKeysCommandOptions {
	options := genKeysCommandOptions{}

	iter := c.args.Iterator()
	for iter.HasNext() {
		nextOpt := iter.Next()
		option, value := args.KeyValue(nextOpt)
		switch option {
		case "--help":
			options.help = true
		case "--privateFile":
			options.privateFile = value
		case "--publicFile":
			options.publicFile = value
		}
	}

	return options
}
