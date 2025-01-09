```bash
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz
```

`~/.zshrc`

```
export PATH=$PATH:/usr/local/go/bin
```

```bash
source ~/.zshrc
```
