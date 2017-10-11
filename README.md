<h1>Copy vault secrets from source to destination.</h1>

Example usage:

Same Vault

```vault-cp --src_token <read-write-token> --src_host <src-host> --src_path /secret/test/web-application --dst_path /secret/production/web-application```

Cross Vault

```vault-cp --src_token <read-token> --dst_token <write-token> --src_host <src-host> --dst_host <dst-host> --src_path /secret/test/web-application --dst_path /secret/production/web-application```


```NAME:
   vault-cp - Copy vault secrets from source to destination
USAGE:
   vault-cp [global options] command [command options] [arguments...]
VERSION:
   0.0.1
COMMANDS:
     help, h  Shows a list of commands or help for one command
GLOBAL OPTIONS:
   --src_token value  Vault token with read access for source [$SRC_VAULT_TOKEN]
   --dst_token value  Vault token with write access for destination, if not specified will use source token [$DST_VAULT_TOKEN]
   --src_host value   Vault source host (default: "http://127.0.0.1:8200")
   --dst_host value   Vault destination host (default: "http://127.0.0.1:8200")
   --src_path value   Source path to copy secrets
   --dst_path value   Destination path to copy secrets
   --help, -h         show help
   --version, -v      print the version
```
