package main

import (
	"fmt"
	"os"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	vault "github.com/hashicorp/vault/api"
	cli "github.com/urfave/cli"
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
			Value:       "",
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

		if destinationVaultToken == "" {
			destinationVaultToken = sourceVaultToken
		}

		if destinationAddr == "" {
			destinationAddr = sourceAddr
		}

		var sourceKeys []string
		var destinationKeys []string
		ch := make(chan string)

		srcClient, dstClient := getClients(sourceAddr, sourceVaultToken, destinationAddr, destinationVaultToken)

		secrets, err := srcClient.Logical().List(sourcePath)
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
			go copySecret(sourceKeys[i], destinationKeys[i], *srcClient, *dstClient, ch)
		}

		for i := 0; i < len(sourceKeys); i++ {
			fmt.Println(<-ch)
		}

		return nil
	}

	app.Run(os.Args)
}

func getClients(sourceAddr string, sourceVaultToken string, destinationAddr string, destinationVaultToken string) (*vault.Client, *vault.Client) {

	srcClient, err := vault.NewClient(&vault.Config{Address: sourceAddr, HttpClient: cleanhttp.DefaultClient()})
	if err != nil {
		fmt.Printf("Could not connect to vault at %v", sourceAddr)
		os.Exit(1)
	}
	srcClient.SetToken(sourceVaultToken)

	dstClient, err := vault.NewClient(&vault.Config{Address: destinationAddr, HttpClient: cleanhttp.DefaultClient()})
	if err != nil {
		fmt.Printf("Could not connect to vault at %v", destinationAddr)
		os.Exit(1)
	}
	dstClient.SetToken(destinationVaultToken)

	return srcClient, dstClient

}
func copySecret(sourceKey string, destinationKey string, srcClient vault.Client, dstClient vault.Client, ch chan string) {

	//Read secret
	secretValue, err := srcClient.Logical().Read(sourceKey)
	if err != nil {
		ch <- fmt.Sprintf("Could not read secert %v %v", sourceKey, err)
		return
	}

	//Write secret
	_, err = dstClient.Logical().Write(destinationKey, secretValue.Data)
	if err != nil {
		ch <- fmt.Sprintf("Could not write secert to  %v %v", destinationKey, err)
		return
	}
	ch <- fmt.Sprintf("Copied secret from %v to %v %v", sourceKey, destinationKey, err)

}
