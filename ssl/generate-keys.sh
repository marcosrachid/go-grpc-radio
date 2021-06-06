
#!/bin/bash
openssl req -x509 -newkey rsa:4096 -keyout private.key \
 -out cert.pem -days 36500 -nodes -subj "/CN=localhost" \
 -extensions EXT -config <( \
  printf "[dn]\nCN=localhost\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:localhost\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")