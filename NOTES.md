NOTES
=====

## Vault

```bash
$ docker run -d --name vault --restart unless-stopped --cap-add=IPC_LOCK -e 'VAULT_LOCAL_CONFIG={"backend": {"file": {"path": "/vault/file"}}, "listener": {"tcp": {"tls_disable": "1"}}, "default_lease_ttl": "168h", "max_lease_ttl": "720h"}' -p 8200:8200 vault server
```

```bash
$ docker exec -it vault sh -c 'VAULT_ADDR="http://127.0.0.1:8200" vault init'
```

```bash
Unseal Key 1: NRgldX5IbREW2FJBOMqYsI7HX3d99yUMmCBg2ut0Hw8B
Unseal Key 2: 3s7eiJOFlq/Jni9GR/+d1bVziFaj29FXPFYt/8rprt4C
Unseal Key 3: ONXS0PZmcJSkMzHS2YNyt7aWwJY0POyQPg41+FCrOooD
Unseal Key 4: jENNKAMWtNq/Ju1eStYE/QW2Dp2T5WvMc2pdKXrkmf4E
Unseal Key 5: alhBcGb1UuHSi/PK1KrrnwZTRl0EAlYLcTJFLuCmDaoF
Initial Root Token: dfd14657-2d94-09a4-649d-64841e28ef36

Vault initialized with 5 keys and a key threshold of 3. Please
securely distribute the above keys. When the Vault is re-sealed,
restarted, or stopped, you must provide at least 3 of these keys
to unseal it again.

Vault does not store the master key. Without at least 3 keys,
your Vault will remain permanently sealed.
```

```bash
$ docker exec -it vault sh -c 'VAULT_ADDR="http://127.0.0.1:8200" vault status'
```

```bash
Sealed: true
Key Shares: 5
Key Threshold: 3
Unseal Progress: 0
Unseal Nonce:
Version: 0.7.0

High-Availability Enabled: true
	Mode: sealed
```