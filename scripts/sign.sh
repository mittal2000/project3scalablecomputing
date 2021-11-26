#!/usr/bin/env sh

sign() {
    openssl genrsa -out $1.key 2048
    openssl req -new -out $1.csr -key $1.key \
        -reqexts SAN \
        -subj "/C=IE/ST=Dublin/L=Dublin/O=TCD/OU=SCSS/CN=$3"
    openssl ca -in $1.csr -out $1.crt \
        -extensions SAN
        
    rm $1.csr
}

main() {
    rm -r $1
    mkdir $1
    sign $1/external.server $2 rasp-*.scss.tcd.ie
    sign $1/internal.server $2 127.*.*.*
    sign $1/client $2 CLIENT
}

main $1 $2