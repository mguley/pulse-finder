#### Steps for using the `key_pair` module

---
##### 1. Generate RSA Key Pair

You need to create RSA key pair (public and private keys) inside the terraform container.
Use the following command to generate a new RSA key pair:

```bash
ssh-keygen -t rsa -b 2048 -f ~/.ssh/pulse-finder-key-pair
```

This command will:
- Generate a private key file: `~/.ssh/pulse-finder-key-pair`
- Generate a public key file: `~/.ssh/pulse-finder-key-pair.pub`

Ensure private key has appropriate file permissions

```bash
chmod 600 ~/.ssh/pulse-finder-key-pair
```

---
##### 2. Use the Public Key in the `key_pair` Module

The `key_pair` module assumes that you already have the public key.
When you use the module:
- Provide the path to the public key file (e.g., `~/.ssh/pulse-finder-key-pair.pub`)
- Terraform will read the public key content and upload it to AWS to create the key pair.

---
##### 3. When connecting to an instance

```bash
ssh -i ~/.ssh/pulse-finder-key-pair ubuntu@<instance-ip>
```
