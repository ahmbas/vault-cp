package main

import (
	"fmt"
	"net/http"
	"os"

	vault "github.com/hashicorp/vault/api"
	"github.com/urfave/cli"
)

func main() {

	var sourceVaultToken string
	var destinationVaultToken string
	var sourceAddr string
	var destinationAddr string
	var sourcePath string
	var destinationPath string

	app := cli.NewApp()
	app.Name = "vault-cp"
	app.Version = "0.0.1"
	app.Usage = "Copy vault secrets from source to destination"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "src_token",
			Value:       "",
			Usage:       "Vault token with read access for source",
			EnvVar:      "SRC_VAULT_TOKEN",
			Destination: &sourceVaultToken,
		},
		cli.StringFlag{
			Name:        "dst_token",
			Value:       sourceVaultToken,
			Usage:       "Vault token with write access for destination, if not specified will use source token",
			EnvVar:      "DST_VAULT_TOKEN",
			Destination: &destinationVaultToken,
		},
		cli.StringFlag{
			Name:        "src_host",
			Value:       "http://127.0.0.1:8200",
			Usage:       "Vault source host",
			Destination: &sourceAddr,
		},
		cli.StringFlag{
			Name:        "dst_host",
			Value:       "http://127.0.0.1:8200",
			Usage:       "Vault destination host",
			Destination: &destinationAddr,
		},
		cli.StringFlag{
			Name:        "src_path",
			Value:       "",
			Usage:       "Source path to copy secrets",
			Destination: &sourcePath,
		},
		cli.StringFlag{
			Name:        "dst_path",
			Value:       "",
			Usage:       "Destination path to copy secrets",
			Destination: &destinationPath,
		},
	}
	app.Action = func(c *cli.Context) error {

		var sourceKeys []string
		var destinationKeys []string
		ch := make(chan string)

		client, err := vault.NewClient(&vault.Config{Address: sourceAddr, HttpClient: http.DefaultClient})
		if err != nil {
			fmt.Printf("Could not connect to vault at %v", sourceAddr)
			os.Exit(1)
		}
		client.SetToken(sourceVaultToken)

		secrets, err := client.Logical().List(sourcePath)
		if err != nil {
			fmt.Printf("Could not fetch secrets from path %v", sourcePath)
			os.Exit(1)
		}

		keys, _ := secrets.Data["keys"]
		list, _ := keys.([]interface{})

		for _, v := range list {
			sourceKeys = append(sourceKeys, sourcePath+"/"+v.(string))
			destinationKeys = append(destinationKeys, destinationPath+"/"+v.(string))
		}

		for i := range sourceKeys {
			go copySecret(sourceKeys[i], destinationKeys[i], *client, ch)
		}

		for i := 0; i < len(sourceKeys); i++ {
			fmt.Println(<-ch)
		}

		return nil
	}

	app.Run(os.Args)
}

func copySecret(sourceKey string, destinationKey string, client vault.Client, ch chan string) {

	//Read secret
	secretValue, err := client.Logical().Read(sourceKey)
	if err != nil {
		ch <- fmt.Sprintf("Could not read secert %v", sourceKey)
	}

	//Write secret
	_, err = client.Logical().Write(destinationKey, secretValue.Data)
	if err != nil {
		ch <- fmt.Sprintf("Could not write secert to  %v", destinationKey)
	}
	ch <- fmt.Sprintf("Copied secret from %v to %v", sourceKey, destinationKey)

}
