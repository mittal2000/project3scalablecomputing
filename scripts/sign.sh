#!/usr/bin/env sh

sign() {
    openssl genrsa -out $1.key 1024
    openssl req -new -out $1.csr -key $1.key \
        -subj "/C=IE/ST=Dublin/L=Dublin/O=TCD/OU=SCSS/CN=$3/"
    openssl x509 -req -in $1.csr -out $1.crt \
        -signkey $1.key -CA $2.crt -CAkey $2.key \
        -extensions v3_req \
        -days 3650 -CAcreateserial
    rm $1.csr
}

main() {
    rm -r $1
    mkdir $1
    sign $1/external.server $2 SERVER
    sign $1/internal.server $2 SERVER
    sign $1/client $2 CLIENT
    cp $2.crt $1
}

main $1 $2