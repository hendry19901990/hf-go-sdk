#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${P2PORT}/$3/" \
        -e "s/\${P3PORT}/$3/" \
        -e "s/\${P4PORT}/$3/" \
        -e "s/\${P5PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $5)
    local CP=$(one_line_pem $6)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${P1PORT}/$3/" \
        -e "s/\${P2PORT}/$3/" \
        -e "s/\${P3PORT}/$3/" \
        -e "s/\${P4PORT}/$3/" \
        -e "s/\${P5PORT}/$3/" \
        -e "s/\${CAPORT}/$4/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ccp-template.yaml | sed -e $'s/\\\\n/\\\n        /g'
}

ORG=1
P0PORT=7051
P1PORT=8051
P2PORT=9051
P3PORT=10051
P4PORT=11051
P5PORT=12051
CAPORT=7054
PEERPEM=crypto-config/peerOrganizations/org1.hf.abl.io/tlsca/tlsca.org1.hf.abl.io-cert.pem
CAPEM=crypto-config/peerOrganizations/org1.hf.abl.io/ca/ca.org1.hf.abl.io-cert.pem

echo "$(json_ccp $ORG $P0PORT $P1PORT $P2PORT $P3PORT $P4PORT $P5PORT $CAPORT $PEERPEM $CAPEM)" > connection-org1.json
echo "$(yaml_ccp $ORG $P0PORT $P1PORT $P2PORT $P3PORT $P4PORT $P5PORT $P2PORT $CAPORT $PEERPEM $CAPEM)" > connection-org1.yaml
